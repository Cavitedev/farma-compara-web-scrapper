package okfarma

import (
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/firestore"
	. "github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/gocolly/colly/v2"
)

const Domain string = "okfarma.es"

var page int = 0

func Scrap(ref *firestore.CollectionRef) {

	fmt.Println(Domain)

	items := []Item{}
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)
	var wg sync.WaitGroup

	c.OnHTML("#product_list", func(h *colly.HTMLElement) {
		fmt.Println("Product List")

		h.ForEach(".product-container", func(_ int, e *colly.HTMLElement) {
			item := Item{}
			item.Name = e.ChildText(".product-name")
			item.Price = e.ChildText(".price")
			item.Url = e.ChildAttr(".product-image-container a", "href")
			wg.Add(1)
			scrapDetailsPage(item.Url, &item, &wg)
			items = append(items, item)
		})
	})

	url := buildPageUrl()
	c.Visit(url)
	wg.Wait()

	bytes, _ := json.Marshal(items)
	fmt.Printf("%+v\n", string(bytes))

}

func scrapDetailsPage(url string, item *Item, wg *sync.WaitGroup) {
	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)

	})

	c.OnHTML("div #center_column", func(h *colly.HTMLElement) {

		item.Image = h.ChildAttr("#bigpic", "src")
		wg.Done()
	})

	c.Visit(url)
}

func buildPageUrl() string {
	page++
	url := fmt.Sprintf("https://%v/medicamentos#/page-%d", Domain, page)
	return url
}
