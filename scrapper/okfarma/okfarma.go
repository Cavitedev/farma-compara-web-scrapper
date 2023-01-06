package okfarma

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	. "github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/gocolly/colly/v2"
)

const Domain string = "okfarma.es"

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
			pageItem := PageItem{}
			pageItem.Website = Domain
			pageItem.Url = e.ChildAttr(".product-image-container a", "href")
			wg.Add(1)
			scrapDetailsPage(&item, &pageItem, &wg)
			item.PageItem = append(item.PageItem, pageItem)
			items = append(items, item)
			time.Sleep(50 * time.Millisecond)
		})
	})

	url := buildPageUrl()
	c.Visit(url)
	wg.Wait()

	bytes, _ := json.Marshal(items)
	fmt.Printf("%+v\n", string(bytes))

}

func scrapDetailsPage(item *Item, pageItem *PageItem, wg *sync.WaitGroup) {
	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)

	})

	c.OnHTML("div #center_column", func(h *colly.HTMLElement) {
		currentTime := time.Now()
		pageItem.LastUpdate = currentTime
		pageItem.Image = h.ChildAttr("#bigpic", "src")
		pageItem.Name = h.ChildText("h1.product-name")
		pageItem.Price = h.ChildText("#our_price_display")
		pageItem.Available = h.ChildText("#availability_value span") != "Este producto ya no está disponible"
		item.Ref = h.ChildAttr("#product_reference>span", "content")
		wg.Done()
	})

	c.Visit(pageItem.Url)
}

func buildPageUrl() string {

	url := fmt.Sprintf("https://%v/medicamentos?id_category=3&n=1192", Domain)
	return url
}
