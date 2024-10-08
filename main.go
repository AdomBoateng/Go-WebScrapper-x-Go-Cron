package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly/v2"
)

type Stock struct {
	company, price, change, change_value string
}

func webscrapper() {
	ticker := []string{
		"PFE",
		"ALTM",
		"NQ=F",
		"GEV",
		"1211.HK",
		"ES=F",
		"MSTR",
		"BTC-USD",
		"ASML",
		"^VIX",
		"EH",
		"LGMK",
		"AMD",
		"LAC",
		"LYFT",
		"RIO",
		"IBRX",
		"BENF",
		"AAPL",
		"OKLO",
		"ASML.AS",
		"PEP",
		"SHOP",
		"IVZ",
		"ADTX",
		"CLSK",
		"TIGR",
		"RACE",
		"WIMI",
		"GBP/USD",
		"USD/JPY",
		"Bitcoin USD",
		"XRP USD",
		"FTSE 100",
		"Nikkei 225",
		"Silver",
		"Gold",
		"VIX",
		"10-Yr Bond",
		"EUR/USD",
		"Crude Oil",
		"Russell 2000 Futures",
		"Nasdaq Futures",
		"Dow Futures",
		"S&P Futures",
	}

	// Create a slice to hold the stock data
	stocks := []Stock{}

	// Create a Colly instance
	c := colly.NewCollector()

	// Handle request logging
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	// Error handling
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error:", err)
	})

	// Parse HTML to extract stock data
	c.OnHTML("section.container.yf-1s1umie", func(e *colly.HTMLElement) {
		stock := Stock{}
		stock.company = e.ChildText("h1.yf-xxbei9")
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")
		stock.change_value = e.ChildText("fin-streamer[data-field='regularMarketChange']")

		// Only append if stock data is found
		if stock.company != "" && stock.price != "" {
			stocks = append(stocks, stock)
		} else {
			fmt.Println("Failed to extract stock data")
		}
	})

	// Visit all tickers
	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}

	// Wait for all requests to complete
	c.Wait()

	// Write to CSV
	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Failed to create CSV file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	headers := []string{"Company", "Price", "Change", "Change-Value"}
	writer.Write(headers)

	for _, stock := range stocks {
		record := []string{stock.company, stock.price, stock.change, stock.change_value}
		writer.Write(record)
	}
	writer.Flush()
}

func main() {
	scheduler := gocron.NewScheduler(time.UTC)

	// Run the web scraper every minute
	scheduler.Every(1).Minute().Do(webscrapper)

	// Start the scheduler
	scheduler.StartBlocking()
}
