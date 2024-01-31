package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type book struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image_url"`
}

func main() {
	c := colly.NewCollector()

	var books []book

	c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
		title := e.ChildText("h3")
		price := e.ChildText(".price_color")
		image := e.ChildAttr("img", "src")

		bookItem := book{
			Title: title,
			Price: price,
			Image: image,
		}
		books = append(books, bookItem)
	})
	
	c.OnHTML("li.next a", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(next_page)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://books.toscrape.com/")
	if err != nil {
		log.Fatal(err)
	}

	// Print the details of the third book if available
	
	content, err := json.Marshal(books)

	if err != nil {
		log.Fatal(err)
	}	

	os.WriteFile("books.json", content, 0644)

}
