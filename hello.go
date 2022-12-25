package main

import (
	"encoding/json"
	"fmt"
	"strings"

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
		// colly.Async(true),
		colly.AllowedDomains(domain),
	)

	fmt.Println("GO")

	items := []item{}

	c.OnHTML("div[class=item-i]", func(h *colly.HTMLElement) {
		var link string = "https://" + domain + h.ChildAttr("a[class=single-image]", "href")
		scrapDetailsPage(link, &items)
	})

	c.OnHTML("div[class=vm-product-details-container]", func(h *colly.HTMLElement) {
		var newItem item = item{}
		// newItem.Price = h.ChildText("span[class=PricesalesPrice]")
		// newItem.Name = h.ChildText("h2[class=product-title]")
		// newItem.ImgUrl = "https://" + domain + h.ChildAttr("img[class=browseProductImage]", "src")
		newItem.Link = h.Request.URL.String()
		newItem.Disponibilidad = h.ChildText("p[class=in-stock].span") == "En stock"
	})

	c.Visit("https://" + domain)
	bytes, _ := json.Marshal(items)
	fmt.Printf("%+v\n", string(bytes))
}
func scrapDetailsPage(url string, items *[]item) {
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)

	})

	c.OnHTML("div .vm-product-container", func(h *colly.HTMLElement) {
		var item item = item{}
		item.Name = h.ChildText("h1")
		item.Price = h.ChildText("span[class=PricesalesPrice]")
		item.Link = h.Request.URL.String()
		item.ImgUrl = domain + h.ChildAttr("#zoom-image", "src")
		var stock string = h.ChildText("p[class=in-stock]")
		item.Disponibilidad = strings.Contains(stock, "En Stock")
		*items = append(*items, item)
	})

	c.Visit(url)
}
