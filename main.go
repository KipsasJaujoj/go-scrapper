package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {

	csvfile, err := os.Create("test.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	w := csv.NewWriter(csvfile)
	c := colly.NewCollector()
	rows := [][]string{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("table[id=main_table_countries_today]", func(table *colly.HTMLElement) {
		table.ForEach("thead", func(_ int, thead *colly.HTMLElement) {
			thead.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
				data := []string{}
				tr.ForEach("th", func(_ int, th *colly.HTMLElement) {
					data = append(data, th.Text)
				})
				rows = append(rows, data)
			})
		})

		table.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			data := []string{}
			tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
				data = append(data, td.Text)
			})
			rows = append(rows, data)
		})

		for _, row := range rows {
			_ = w.Write(row)
		}
		w.Flush()
		csvfile.Close()
	})
	c.Visit("https://www.worldometers.info/coronavirus/#countries")
}
