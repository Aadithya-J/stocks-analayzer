package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change, changePercent string
}

func main() {

	tickers := []string{
		"MSFT", "IBM", "AAPL", "GOOGL", "AMZN",
	}
	// 	"FB", "NFLX", "TSLA", "NVDA", "INTC",
	// 	"AMD", "QCOM", "ADBE", "PYPL", "CRM",
	// 	"ORCL", "V", "MA", "ACN", "CTSH",
	// 	"INFY", "ADP", "PAYX", "CDNS", "ZOMATO.NS",
	// }

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		fmt.Println("Visiting:", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	stocks := []Stock{}
	c.OnHTML(".gridLayout .container", func(e *colly.HTMLElement) {
		stock := Stock{}
		e.DOM.Find(".top").Each(func(i int, top *goquery.Selection) {
			top.Find(".hdr .left .container h1").Each(func(i int, h1 *goquery.Selection) {
				company := h1.Text()
				fmt.Println("\nCompany Data:")
				fmt.Println(company)
				stock.company = company
			})
		})
		e.DOM.Find(".bottom").Each(func(i int, bottom *goquery.Selection) {
			bottom.Find(`[data-testid="qsp-price"]`).Each(func(i int, s *goquery.Selection) {
				fmt.Println("\nPrice Data:")
				fmt.Println(s.Text())
				stock.price = s.Text()
			})
			bottom.Find(`[data-testid="qsp-price-change"]`).Each(func(i int, s *goquery.Selection) {
				fmt.Println("\nPrice Change Data:")
				fmt.Println(s.Text())
				stock.change = s.Text()
			})

			bottom.Find(`[data-testid="qsp-price-change-percent"]`).Each(func(i int, s *goquery.Selection) {
				fmt.Println("\nPrice Change Percent Data:")
				fmt.Println(s.Text())
				stock.changePercent = s.Text()
			})
		})
		if stock.company != "" && stock.price != "" {
			stocks = append(stocks, stock)
		}
	})

	for _, ticker := range tickers {
		url := "https://finance.yahoo.com/quote/" + ticker
		c.Visit(url)
		time.Sleep(2 * time.Second)
	}
	c.Wait()

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, stock := range stocks {
		row := []string{stock.company, stock.price, stock.change, stock.changePercent}
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	fmt.Println("Data has been written to stocks.csv")
}
