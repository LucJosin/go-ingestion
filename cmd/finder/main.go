package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/LucJosin/go-ingestion/internal/model"
	"github.com/gocolly/colly/v2"
)

func main() {
	profileUrl := flag.String("p", "", "Profile to find")
	flag.Parse()

	var bank model.Bank
	if *profileUrl != "" {
		bank.Profile = *profileUrl
	} else {
		if err := json.NewDecoder(os.Stdin).Decode(&bank); err != nil {
			log.Fatalf("Failed to parse input: %v", err)
		}
	}

	c := colly.NewCollector()

	profile := model.BankProfile{
		Name:  bank.Name,
		Lists: []model.BankListData{},
	}

	// bank name
	c.OnHTML("h1.listuser-header__name", func(e *colly.HTMLElement) {
		profile.Name = e.Text
	})

	// employees and ceo
	c.OnHTML("dl.listuser-block__item", func(e *colly.HTMLElement) {
		profileStatsTitle := e.ChildText(".profile-stats__title")
		profileStatsText := e.ChildText(".profile-stats__text span.profile-stats__text")

		if profileStatsTitle == "Employees" {
			profile.Employees, _ = strconv.Atoi(profileStatsText)
		}

		if profileStatsTitle == "CEO" {
			profile.CEO = profileStatsText
		}
	})

	// tagged list
	c.OnHTML("div.ranking .listuser-block__item", func(e *colly.HTMLElement) {
		rankingName := e.ChildText(".listuser-item__list--title")
		rankingUrl := e.ChildAttr(".listuser-item__list--title", "href")

		profile.Lists = append(profile.Lists, model.BankListData{
			Name: rankingName,
			URL:  rankingUrl,
		})
	})

	err := c.Visit(bank.Profile)
	if err != nil {
		return
	}

	jsonData, err := json.Marshal(profile)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(jsonData))
}
