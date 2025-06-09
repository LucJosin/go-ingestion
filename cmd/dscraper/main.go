package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/LucJosin/go-ingestion/internal/exporter"
	"github.com/LucJosin/go-ingestion/internal/model"
	"golang.org/x/net/html"
)

type industryRank struct {
	Industry string `json:"industry"`
	Rank     int    `json:"rank"`
}

// simple data > table data format
type content struct {
	OrganizationName string         `json:"organizationName"`
	Uri              string         `json:"uri"`
	City             string         `json:"city"`
	Country          string         `json:"country"`
	YearFounded      int            `json:"yearFounded"`
	IndustryRanks    []industryRank `json:"industryRanks"`
}

func (c content) toDTO() model.Bank {
	var rank = 0
	for _, h := range c.IndustryRanks {
		if h.Industry == c.Country {
			rank = h.Rank
		}
	}

	return model.Bank{
		Name:    c.OrganizationName,
		City:    c.City,
		Country: c.Country,
		Founded: c.YearFounded,
		Rank:    rank,
		Profile: fmt.Sprintf("https://www.forbes.com/companies/%s/?list=worlds-best-banks", c.Uri),
	}
}

// Data Scraper (using the injected content)
//
// For this "dscraper" we can use the json data that is injected into the HTML page,
// you can preview the data using the devtools console with the command:
//
// window["forbes"]["simple-site"].tableData
//
// This command line just parse the HTML, finder for the script and export the data
func main() {
	output := flag.String("o", "stdout", "Output format")
	filename := flag.String("f", "companies", "Output file name")
	flag.Parse()

	expr := exporter.NewExporter(*output, *filename)

	res, err := http.Get("https://www.forbes.com/lists/worlds-best-banks/")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()

	page, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
		return
	}

	var data []content
	searchForScript(page, &data)

	var companies []model.Bank
	for _, d := range data {
		companies = append(companies, d.toDTO())
	}

	err = expr.ExportData(companies)
	if err != nil {
		log.Println(err)
		return
	}
}

// the content in the page are injected into the html, we can just take and parse this data
func searchForScript(n *html.Node, output any) {
	if n.Type == html.ElementNode && n.Data == "script" {
		// window["forbes"]["simple-site"].tableData
		if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "simple-site") {
			pattern := regexp.MustCompile(`"tableData":(\[{.*?}]}])`)
			matches := pattern.FindStringSubmatch(n.FirstChild.Data)

			if len(matches) > 1 {
				jsonData := matches[1]

				err := json.Unmarshal([]byte(jsonData), &output)
				if err != nil {
					log.Println("Error unmarshalling JSON:", err)
				}

				return
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		searchForScript(c, output)
	}
}
