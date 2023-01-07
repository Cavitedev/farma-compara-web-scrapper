package okfarma

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	. "github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

const Domain string = "okfarma.es"

type WaitGroupCount struct {
	sync.WaitGroup
	count int32
}

func Scrap(ref *firestore.CollectionRef) {

	fmt.Println(Domain)

	items := []Item{}
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	c.OnHTML("#product_list", func(h *colly.HTMLElement) {
		fmt.Println("Product List")

		h.ForEach(".product-container", func(_ int, e *colly.HTMLElement) {
			item := Item{}
			pageItem := PageItem{}
			pageItem.Website = Domain
			pageItem.Url = e.ChildAttr(".product-image-container a", "href")
			scrapDetailsPage(&item, &pageItem)
			item.PageItem = append(item.PageItem, pageItem)
			items = append(items, item)
			time.Sleep(50 * time.Millisecond)

		})
	})

	url := buildPageUrl()
	c.Visit(url)

	bytes, _ := json.Marshal(items)
	fmt.Printf("%+v\n", string(bytes))

}

func scrapDetailsPage(item *Item, pageItem *PageItem) {
	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)
	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)

	})

	c.OnHTML("div #center_column", func(h *colly.HTMLElement) {
		currentTime := time.Now()
		pageItem.LastUpdate = currentTime
		pageItem.Image = h.ChildAttr("#bigpic", "src")
		pageItem.Name = h.ChildText("h1.product-name")

		price := h.ChildText("#our_price_display")
		pageItem.Price = utils.SpanishNumberStrToNumber(price)
		pageItem.Available = h.ChildText("#availability_value span") != "Este producto ya no está disponible"
		item.Ref = h.ChildAttr("#product_reference>span", "content")
	})

	c.Visit(pageItem.Url)
}

func buildPageUrl() string {

	url := fmt.Sprintf("https://%v/medicamentos", Domain)
	return url
}
