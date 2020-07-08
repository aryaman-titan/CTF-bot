package main

import (
	"./scrapper"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("mybot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			fmt.Println(parts)
			if len(parts) == 3 && strings.ToLower(parts[1]) == "info" && strings.ToLower(parts[2]) == "upcomingctf" {
				// looks good, get the quote and reply with the result
				go func(m Message) {
					m.Text = jsonCtfFormat(scrapper.GetCTFs())
					fmt.Println(m.Text)
					postMessage(ws, m)
				}(m)
				// NOTE: the Message object is copied, this is intentional
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				postMessage(ws, m)
			}
		}
	}
}

func jsonCtfFormat(byteValue []byte) string {
	var CTFLists []scrapper.Ctf
	json.Unmarshal(byteValue, &CTFLists)

	result := ""
	for i := 0; i < len(CTFLists); i++ {
		result += fmt.Sprintf("*Name*: %s\n*Date*: %s\n*Duration*: %s\n\n", CTFLists[i].Name, CTFLists[i].Date, CTFLists[i].Duration)
		fmt.Println(result)
	}

	return result
}
