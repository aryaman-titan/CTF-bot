package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("!")
}

type ctf struct {
	Name    string
	Date    string
	Duarion string
}

func jsonMarshal(rows [][]string) {
	for _, p := range rows {
		m := ctf{p[0], p[1], p[2]}
		var jsonData []byte
		jsonData, _ = json.Marshal(m)
		fmt.Println(string(jsonData))

	}
}
func main() {
	var row []string
	var rows [][]string
	client := &http.Client{
		CheckRedirect: noRedirect,
	}
	req, _ := http.NewRequest("GET", "https://ctftime.org/", nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, _ := client.Do(req)

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("table.upcoming-events").Each(func(index int, tablehtml *goquery.Selection) {
		if index == 0 {
			tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
				rowhtml.Find("td").Each(func(i int, tableheading *goquery.Selection) {
					if i != 0 {

						temp := tableheading.Text()
						if i == 2 {
							index := strings.Index(temp, "UTC")
							temp1 := temp[:index-1]
							temp = temp1
						}

						row = append(row, temp)
					}
				})
				if row != nil {
					rows = append(rows, row)
					row = nil
				}
			})
		}
	})
	jsonMarshal(rows)
}
