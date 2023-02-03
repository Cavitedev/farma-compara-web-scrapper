package dosfarma

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	. "github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

const websiteName string = "dosfarma"
const Domain string = "www.dosfarma.com"

var itemCount int = 1
var page int = 1

func Scrap(ref *firestore.CollectionRef) {

	log.Println(Domain)

	items := []Item{}
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	c.OnHTML("#js-product-list", func(h *colly.HTMLElement) {

		itemCount = 0

		h.ForEach(".item", func(_ int, e *colly.HTMLElement) {

			item := Item{}
			pageItem := WebsiteItem{}
			pageItem.Url = e.ChildAttr(".product-thumbnail", "href")
			scrapDetailsPage(&item, &pageItem)
			if item.WebsiteItems == nil {
				item.WebsiteItems = make(map[string]WebsiteItem)
			}
			item.WebsiteItems[websiteName] = pageItem
			items = append(items, item)
			firestore_utils.UpdateItem(item, ref)
			time.Sleep(50 * time.Millisecond)
			h.Attr("class")
			itemCount++
		})

		log.Printf("Scrapped %v items", itemCount)

	})

	for itemCount > 0 {
		c.Visit(fmt.Sprintf("https://www.dosfarma.com/higiene/corporal/?page=%v", page))
		page++
		log.Printf("Scrapped %v items on page %v", itemCount, page)
	}

}

var productsVisited int = 0

func scrapDetailsPage(item *Item, pageItem *WebsiteItem) {
	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)
	c.OnResponse(func(r *colly.Response) {
		productsVisited++
		log.Printf("Visit %d URL:%v\n", productsVisited, r.Request.URL)

	})

	c.OnHTML("#add-to-cart-or-refresh", func(h *colly.HTMLElement) {
		currentTime := time.Now()
		pageItem.LastUpdate = currentTime
		pageItem.Image = h.ChildAttr("img", "src")
		pageItem.Name = h.ChildText("h1")

		price := h.ChildText(".final-price")
		pageItem.Price = utils.ParseSpanishNumberStrToNumber(price)
		pageItem.Available = h.ChildTexts(".disponible")[0] == "En stock"
		item.Ref = h.ChildTexts(".referencia")[0]
	})

	c.Visit(pageItem.Url)
}
