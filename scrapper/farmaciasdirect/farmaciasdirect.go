package farmaciasdirect

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

const Domain string = "www.farmaciasdirect.com"

var pageNum int = 1

func Scrap(ref *firestore.CollectionRef) {

	log.Println(Domain)
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	c.OnHTML("#js-product-list", func(h *colly.HTMLElement) {
		log.Println("Product List")

		h.ForEach(".card-product", func(_ int, e *colly.HTMLElement) {
			item := Item{}
			pageItem := WebsiteItem{}
			pageItem.Url = e.ChildAttr(".card-body>a", "href")
			scrapDetailsPage(&item, &pageItem)
			if item.WebsiteItems == nil {
				item.WebsiteItems = make(map[string]WebsiteItem)
			}
			item.WebsiteItems[Domain] = pageItem
			firestore_utils.UpdateItem(item, ref)
			time.Sleep(50 * time.Millisecond)
		})
	})

	url := buildPageUrl()
	c.Visit(url)
}

func scrapDetailsPage(item *Item, pageItem *WebsiteItem) {
	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)
	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)

	})

	c.OnHTML("#main", func(h *colly.HTMLElement) {
		currentTime := time.Now()
		pageItem.LastUpdate = currentTime
		pageItem.Image = h.ChildAttr("img.img-fluid", "src")
		pageItem.Name = h.ChildText("h1.product-name")

		price := h.ChildAttr(".current-price>span", "content")
		pageItem.Price = utils.ParseSpanishNumberStrToNumber(price)
		pageItem.Available = h.ChildText("#product-availability") == ""
		item.Ref = h.ChildTexts("div.product-reference>span")[0]
	})

	c.Visit(pageItem.Url)
}

func buildPageUrl() string {

	url := fmt.Sprintf("https://%v/medicamentos-8?page=%v", Domain, pageNum)
	pageNum++
	return url
}
