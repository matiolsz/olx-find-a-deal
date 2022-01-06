package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	name  string
	price string
}

func main() {
	items := make([]Item, 0)

	url := "https://www.olx.pl/motoryzacja/samochody/q-porsche-911/?search%5Bfilter_float_price%3Afrom%5D=300000&search%5Bfilter_float_price%3Ato%5D=400000"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".container tr .offer-wrapper ").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").Text()
		price := s.Find("p.price").Text()

		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		processedTitle := reg.ReplaceAllString(title, "")
		processedPrice := reg.ReplaceAllString(price, "")
		processedPrice = strings.TrimRight(processedPrice, "z")
		var smallTitle string
		var item Item
		if len(processedTitle) > 40 {
			smallTitle = processedTitle[0:40]
			item = Item{smallTitle, processedPrice}
		} else {
			item = Item{processedTitle, processedPrice}
		}
		items = append(items, item)

	})

	for i, val := range items {
		fmt.Printf("%-3d: name: %-40s\tprice: %s\n", i, val.name, val.price)
	}
}
