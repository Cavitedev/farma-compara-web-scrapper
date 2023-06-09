package farmaciaencasa

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	. "github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

const websiteName string = "farmaciaencasa"
const Domain string = "www.farmaciaencasaonline.es"

var lastPage int = 5
var page int = 1

func Scrap(client *firestore.Client) {

	log.Println(Domain)

	items := []Item{}
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	c.OnResponse(func(r *colly.Response) {
		log.Printf("Visit URL:%v\n", r.Request.URL)

	})

	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML(".pages", func(h *colly.HTMLElement) {
		pagesLi := h.ChildTexts("li>a")
		lastPageLi := pagesLi[len(pagesLi)-2]
		lastPageLi = utils.NumberRegexString(lastPageLi)
		lastPageI64, err := strconv.ParseInt(lastPageLi, 10, 32)
		if err != nil {
			log.Println("Error parsing " + lastPageLi)
		}
		if lastPageI64 > int64(lastPage) {
			lastPage = int(lastPageI64)
		}

	})

	c.OnHTML(".product-items", func(h *colly.HTMLElement) {

		h.ForEach(".product-item", func(_ int, e *colly.HTMLElement) {

			item := Item{}
			pageItem := WebsiteItem{}
			pageItem.Url = e.ChildAttr(".product", "href")
			scrapDetailsPage(&item, &pageItem)
			if item.WebsiteItems == nil {
				item.WebsiteItems = make(map[string]WebsiteItem)
			}
			item.WebsiteItems[websiteName] = pageItem
			items = append(items, item)
			firestore_utils.UpdateItem(item, client)
			time.Sleep(50 * time.Millisecond)
			h.Attr("class")
		})

	})

	for page != lastPage+1 {
		c.Visit(fmt.Sprintf("https://www.farmaciaencasaonline.es/corporal/cuerpo?p=%v", page))
		page++
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

	c.OnHTML(".gallery-placeholder__image", func(h *colly.HTMLElement) {
		pageItem.Image = h.Attr("src")
	})

	c.OnHTML(".product-info-main", func(h *colly.HTMLElement) {

		references := h.ChildTexts(".sku>div")
		if len(references) == 0 {
			return
		}
		item.Ref = references[0]

		currentTime := time.Now()
		pageItem.LastUpdate = currentTime
		pageItem.Name = h.ChildText(".page-title>span")

		price := h.ChildText(".price")
		if price == "" {
			price = h.ChildText(".product-type-data>div>.regular-price>.price")
		}
		pageItem.Price = utils.ParseSpanishNumberStrToNumber(price)

		availableTexts := h.ChildTexts(".stock>span")
		if len(availableTexts) > 0 {
			pageItem.Available = availableTexts[0] == "Disponible"
		}

	})

	c.Visit(pageItem.Url)
}
