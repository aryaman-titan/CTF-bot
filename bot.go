package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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
			fmt.Println(parts[1])
			if len(parts) == 3 && parts[1] == "info" {
				// looks good, get the quote and reply with the result
				go func(m Message) {
					// m.Text = getQuote(parts[2])
					 m.Text = getCTF()
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

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("!")
}

// Users struct which contains
// an array of users
type Users struct {
    Users []User `json:"users"`
}

// User struct which contains a name
// a type and a list of social links
type User struct {
    Name   string `json:"name"`
    Type   string `json:"type"`
    Age    int    `json:"Age"`
    Social Social `json:"social"`
}

// Social struct which contains a
// list of links
type Social struct {
    Facebook string `json:"facebook"`
    Twitter  string `json:"twitter"`
}

func getQuote(sym string) string {
	jsonFile,err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened ",jsonFile.Name())

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users Users

	json.Unmarshal(byteValue, &users)

	n,_ := strconv.Atoi(sym)
	return fmt.Sprintf("Name: %s \nNickname: %s",users.Users[n].Name, users.Users[n].Type)
}

func getCTF() (string) {
	client := &http.Client{
		CheckRedirect: noRedirect,
	}
	req, _ := http.NewRequest("GET", "https://ctftime.org/", nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	var result string
	doc.Find("table.upcoming-events").Each(func(index int, tablehtml *goquery.Selection) {
		if index == 0 {
			result = fmt.Sprintf("Here's the list, %s",tablehtml.Text())
			fmt.Println(result)
		}
	})

	return result

}