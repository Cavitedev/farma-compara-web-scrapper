package main

import (
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly/v2"
)

type item struct {
	Name           string `json:"name"`
	Price          string `json:"price"`
	ImgUrl         string `json:"imgUrl"`
	Link           string `json:"link"`
	Disponibilidad bool   `json:"disponibilidad"`
}

var domain string = "www.farmaciaalbacete.es"

func main() {

	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains(domain),
	)

	fmt.Println("GO")

	items := []item{}

	c.OnHTML("div[class=item-i]", func(h *colly.HTMLElement) {
		var newItem item = item{}
		newItem.Price = h.ChildText("span[class=PricesalesPrice]")
		newItem.Name = h.ChildText("h2[class=product-title]")
		newItem.ImgUrl = "https://" + domain + h.ChildAttr("img[class=browseProductImage]", "src")
		newItem.Link = "https://" + domain + h.ChildAttr("a[class=single-image]", "href")
		c.Visit(newItem.Link)
		items = append(items, newItem)
	})

	c.OnHTML("div[class=vm-product-details-container", func(h *colly.HTMLElement) {
		link := h.Request.URL.String()
		item := items[0]
		item.Disponibilidad = h.ChildText("p[class=in-stock].span") == "En stock"
	})

	c.Visit("https://" + domain)
	bytes, _ := json.Marshal(items)
	fmt.Printf("%+v\n", string(bytes))
}
func scrapDetailsPage(url string, item *item) {
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)
	c.OnHTML("div[class=vm-product-details-container", func(h *colly.HTMLElement) {
		item.Disponibilidad = h.ChildText("p[class=in-stock].span") == "En stock"
	})

	c.Visit(url)
}
