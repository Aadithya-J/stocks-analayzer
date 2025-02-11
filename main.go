package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

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

	var stocks []Stock
	var mu sync.Mutex
	var wg sync.WaitGroup

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error:", err)
	})

	c.OnHTML(".gridLayout .container", func(e *colly.HTMLElement) {
		stock := Stock{}
		stock.company = e.DOM.Find(".hdr .left .container h1").Text()
		stock.price = e.DOM.Find(`[data-testid="qsp-price"]`).Text()
		stock.change = e.DOM.Find(`[data-testid="qsp-price-change"]`).Text()
		stock.changePercent = e.DOM.Find(`[data-testid="qsp-price-change-percent"]`).Text()

		if stock.company != "" && stock.price != "" {
			mu.Lock()
			stocks = append(stocks, stock)
			mu.Unlock()
		}
	})

	for _, ticker := range tickers {
		wg.Add(1)
		go func(ticker string) {
			defer wg.Done()
			c.Visit("https://finance.yahoo.com/quote/" + ticker)
		}(ticker)
	}

	wg.Wait()

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatal("Cannot create file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, stock := range stocks {
		if err := writer.Write([]string{stock.company, stock.price, stock.change, stock.changePercent}); err != nil {
			log.Fatal("Cannot write to file:", err)
		}
	}

	fmt.Println("Data has been written to stocks.csv")
}
