package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/LucJosin/go-ingestion/internal/exporter"
	"github.com/LucJosin/go-ingestion/internal/model"
	"github.com/gocolly/colly/v2"
)

func main() {
	output := flag.String("o", "stdout", "Output format")
	filename := flag.String("f", "companies", "Output file name")
	flag.Parse()

	expr := exporter.NewExporter(*output, *filename)

	c := colly.NewCollector()

	var companies []model.Bank

	c.OnHTML("div.table-row-group", func(e *colly.HTMLElement) {
		e.ForEach("a.table-row", func(_ int, a *colly.HTMLElement) {
			company := model.Bank{}
			company.Name = a.ChildText(".organizationName .row-cell-value")
			company.City = a.ChildText(".city .row-cell-value")
			company.Country = a.ChildText(".country .row-cell-value")

			founded := a.ChildText(".yearFounded .row-cell-value")
			company.Founded, _ = strconv.Atoi(founded)

			// div.industryRankCell
			//  |-> span.starRank
			//  |-> span.starRankIndustry
			a.ForEach(".industryRankCell", func(_ int, i *colly.HTMLElement) {
				company.Country = i.ChildText(".starRankIndustry")

				if rank := i.ChildText(".starRank"); rank != "" {
					company.Rank, _ = strconv.Atoi(rank)
				}
			})

			// some rows don't contain the 'href' value
			company.Profile = a.Attr("href")
			if company.Profile == "" {
				uri := a.Attr("uri")
				company.Profile = fmt.Sprintf("https://www.forbes.com/companies/%s/?list=worlds-best-banks", uri)
			}

			companies = append(companies, company)
		})
	})

	err := c.Visit("https://www.forbes.com/lists/worlds-best-banks/")
	if err != nil {
		return
	}

	err = expr.ExportData(companies)
	if err != nil {
		log.Println(err)
		return
	}
}
