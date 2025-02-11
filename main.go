package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

func main1() {
	ticker := "ZOMATO.NS"
	url := "https://finance.yahoo.com/quote/" + ticker

	c := colly.NewCollector()
	var rawHTML string

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("Visiting", req.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		rawHTML = string(r.Body)

		err := os.WriteFile("output.html", r.Body, 0644)
		if err != nil {
			log.Println("Error writing to file:", err)
		}
	})

	c.Visit(url)
	c.Wait()

	priceRegex := regexp.MustCompile(`data-testid="qsp-price">([\d.,]+)`)
	changeRegex := regexp.MustCompile(`data-testid="qsp-price-change">(-?[\d.,]+)`)
	changePercentRegex := regexp.MustCompile(`data-testid="qsp-price-change-percent">\(([-+]?[\d.,]+%)\)`)

	priceMatch := priceRegex.FindStringSubmatch(rawHTML)
	price := "N/A"
	if len(priceMatch) > 1 {
		price = priceMatch[1]
	}

	changeMatch := changeRegex.FindStringSubmatch(rawHTML)
	change := "N/A"
	if len(changeMatch) > 1 {
		change = changeMatch[1]
	}

	changePercentMatch := changePercentRegex.FindStringSubmatch(rawHTML)
	changePercent := "N/A"
	if len(changePercentMatch) > 1 {
		changePercent = changePercentMatch[1]
	}

	fmt.Printf("Stock: %s\nPrice: %s\nChange: %s\nChange Percent: %s\n", ticker, price, change, changePercent)
}
