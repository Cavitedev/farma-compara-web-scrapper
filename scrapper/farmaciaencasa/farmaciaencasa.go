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

func Scrap(ref *firestore.CollectionRef) {

	log.Println(Domain)

	items := []Item{}
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	c.OnResponse(func(r *colly.Response) {
		log.Printf("Visit URL:%v\n", r.Request.URL)

	})

	c.OnHTML(".pages", func(h *colly.HTMLElement) {
		pagesLi := h.ChildTexts("li>a")
		lastPageLi := pagesLi[len(pagesLi)-2]
		lastPageI64, err := strconv.ParseInt(lastPageLi, 10, 32)
		if err != nil {
			log.Println("Error parsing " + lastPageLi)
		}
		lastPage = int(lastPageI64)

	})

	c.OnHTML(".itemgrid", func(h *colly.HTMLElement) {

		h.ForEach(".item", func(_ int, e *colly.HTMLElement) {

			item := Item{}
			pageItem := WebsiteItem{}
			pageItem.Url = e.ChildAttr(".product-name>a", "href")
			scrapDetailsPage(&item, &pageItem)
			if item.WebsiteItems == nil {
				item.WebsiteItems = make(map[string]WebsiteItem)
			}
			item.WebsiteItems[websiteName] = pageItem
			items = append(items, item)
			firestore_utils.UpdateItem(item, ref)
			time.Sleep(50 * time.Millisecond)
			h.Attr("class")
		})

	})

	for page != lastPage {
		c.Visit(fmt.Sprintf("https://www.farmaciaencasaonline.es/corporal-cuidado-cuerpo-c-103_126.html?limit=60&p=%v", page))
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

	c.OnHTML(".product-view", func(h *colly.HTMLElement) {

		references := h.ChildTexts(".sku>span>span")
		if len(references) == 0 {
			return
		}
		item.Ref = references[0]

		currentTime := time.Now()
		pageItem.LastUpdate = currentTime
		pageItem.Image = h.ChildAttr("img", "src")
		pageItem.Name = h.ChildText(".product-name>h1")

		price := h.ChildText(".special-price>.price")
		if price == "" {
			price = h.ChildText(".product-type-data>div>.regular-price>.price")
		}
		pageItem.Price = utils.ParseSpanishNumberStrToNumber(price)

		availableTexts := h.ChildTexts(".availability>span")
		if len(availableTexts) > 0 {
			pageItem.Available = availableTexts[0] == "En existencia"
		}

	})

	c.Visit(pageItem.Url)
}
