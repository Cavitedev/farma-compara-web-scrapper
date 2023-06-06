package okfarma

import (
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	. "github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

const websiteName string = "okfarma"
const Domain string = "okfarma.es"

func Scrap(ref *firestore.CollectionRef) {

	log.Println(Domain)

	items := []Item{}
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	c.SetRequestTimeout(100 * time.Second)

	c.OnHTML("a.product-name", func(h *colly.HTMLElement) {
		item := Item{}
		pageItem := WebsiteItem{}
		pageItem.Url = h.Attr("href")
		scrapDetailsPage(&item, &pageItem)
		if item.WebsiteItems == nil {
			item.WebsiteItems = make(map[string]WebsiteItem)
		}
		item.WebsiteItems[websiteName] = pageItem
		items = append(items, item)
		firestore_utils.UpdateItem(item, ref)
		time.Sleep(50 * time.Millisecond)

		log.Printf("start \n")
		h.Attr("class")

	})

	c.Visit("https://okfarma.es/higiene-corporal?id_category=14&n=10000")
	log.Printf("Scrapped %v items", len(items))

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

	c.OnHTML("div #center_column", func(h *colly.HTMLElement) {
		currentTime := time.Now()
		pageItem.LastUpdate = currentTime
		pageItem.Image = h.ChildAttr("#bigpic", "src")
		pageItem.Name = h.ChildText("h1.product-name")

		price := h.ChildText("#our_price_display")
		pageItem.Price = utils.ParseSpanishNumberStrToNumber(price)
		pageItem.Available = h.ChildText("#availability_value span") != "Este producto ya no está disponible"
		item.Ref = h.ChildAttr("#product_reference>span", "content")
	})

	c.Visit(pageItem.Url)
}
