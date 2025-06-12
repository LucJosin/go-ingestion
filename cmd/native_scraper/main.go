package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/LucJosin/go-ingestion/internal/exporter"
	"github.com/LucJosin/go-ingestion/internal/model"
	"golang.org/x/net/html"
)

var banks []model.Bank

func main() {
	output := flag.String("o", "stdout", "Output format")
	filename := flag.String("f", "banks", "Output file name")
	flag.Parse()

	expr := exporter.NewExporter(*output, *filename)

	res, err := http.Get("https://www.forbes.com/lists/worlds-best-banks/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	page, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	runScraper(page)

	if err = expr.ExportData(banks); err != nil {
		log.Fatal(err)
	}
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, c := range classes {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

func runScraper(n *html.Node) {
	if t := findTableElement(n); t == nil {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			runScraper(c)
		}
	}
}

func findTableElement(n *html.Node) *html.Node {
	if hasClass(n, "table") {
		for t := n.FirstChild; t != nil; t = t.NextSibling {
			findTableRowGroup(t)
		}
	}
	return nil
}

func findTableRowGroup(n *html.Node) *html.Node {
	if hasClass(n, "table-row-group") {
		for trg := n.FirstChild; trg != nil; trg = trg.NextSibling {
			findTableRow(trg)
		}
	}
	return nil
}

func findTableRow(n *html.Node) *html.Node {
	if hasClass(n, "table-row") {
		var bank model.Bank

		for _, attr := range n.Attr {
			if attr.Key == "href" && attr.Val != "" {
				bank.Profile = attr.Val
				break
			}

			if attr.Key == "uri" && attr.Val != "" {
				bank.Profile = fmt.Sprintf("https://www.forbes.com/companies/%s/?list=worlds-best-banks", attr.Val)
				break
			}
		}

		for tr := n.FirstChild; tr != nil; tr = tr.NextSibling {
			extractData(tr, &bank)
		}

		banks = append(banks, bank)
	}
	return nil
}

func extractData(n *html.Node, bank *model.Bank) {
	rcv := func(rn *html.Node) (string, error) {
		for ns := rn.FirstChild; ns != nil; ns = ns.NextSibling {
			if hasClass(ns, "row-cell-value") {
				return ns.FirstChild.Data, nil
			}
		}

		return "", errors.New("no row-cell-value found")
	}

	if hasClass(n, "organizationName") {
		if value, err := rcv(n); err == nil {
			bank.Name = value
		}
	}

	if hasClass(n, "city") {
		if value, err := rcv(n); err == nil {
			bank.City = value
		}
	}

	if hasClass(n, "country") {
		if value, err := rcv(n); err == nil {
			bank.Country = value
		}
	}

	if hasClass(n, "yearFounded") {
		if value, err := rcv(n); err == nil {
			bank.Founded, _ = strconv.Atoi(value)
		}
	}

	if hasClass(n, "industryRanks") {
		for ir := n.FirstChild; ir != nil; ir = ir.NextSibling {
			if hasClass(ir, "industryRankCell") {
				var industry string
				for irc := ir.FirstChild; irc != nil; irc = irc.NextSibling {
					if hasClass(irc, "starRankIndustry") {
						industry = irc.FirstChild.Data
						continue
					}

					if hasClass(irc, "starRank") && industry == bank.Country {
						bank.Rank, _ = strconv.Atoi(irc.FirstChild.Data)
						break
					}
				}
			}
		}
	}
}
