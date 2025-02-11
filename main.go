package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change, changePercent string
}

func main() {
	tickers := []string{
		"MSFT", "IBM", "AAPL", "GOOGL", "AMZN",
		"FB", "NFLX", "TSLA", "NVDA", "INTC",
		"AMD", "QCOM", "ADBE", "PYPL", "CRM",
		"ORCL", "V", "MA", "ACN", "CTSH",
		"INFY", "ADP", "PAYX", "CDNS", "ZOMATO.NS",
	}

	c := colly.NewCollector()
	rawHTMLs := []string{}

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("Visiting", req.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		rawHTML := string(r.Body)
		rawHTMLs = append(rawHTMLs, rawHTML)
	})

	for _, ticker := range tickers {
		url := "https://finance.yahoo.com/quote/" + ticker
		c.Visit(url)
	}
	c.Wait()

	priceRegex := regexp.MustCompile(`data-testid="qsp-price">([\d.,]+)`)
	changeRegex := regexp.MustCompile(`data-testid="qsp-price-change">(-?[\d.,]+)`)
	changePercentRegex := regexp.MustCompile(`data-testid="qsp-price-change-percent">\(([-+]?[\d.,]+%)\)`)

	stocks := []Stock{}
	for _, rawHTML := range rawHTMLs {
		stock := Stock{}
		priceMatch := priceRegex.FindStringSubmatch(rawHTML)
		stock.price = "N/A"
		if len(priceMatch) > 1 {
			stock.price = priceMatch[1]
		}

		changeMatch := changeRegex.FindStringSubmatch(rawHTML)
		stock.change = "N/A"
		if len(changeMatch) > 1 {
			stock.change = changeMatch[1]
		}

		changePercentMatch := changePercentRegex.FindStringSubmatch(rawHTML)
		stock.changePercent = "N/A"
		if len(changePercentMatch) > 1 {
			stock.changePercent = changePercentMatch[1]
		}

		stocks = append(stocks, stock)
	}

	for i, stock := range stocks {
		fmt.Printf("Stock %d:\nPrice: %s\nChange: %s\nChange Percent: %s\n\n", i+1, stock.price, stock.change, stock.changePercent)
	}

	fmt.Println(stocks)

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, stock := range stocks {
		row := []string{stock.price, stock.change, stock.changePercent}
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	fmt.Println("Data has been written to stocks.csv")
}
