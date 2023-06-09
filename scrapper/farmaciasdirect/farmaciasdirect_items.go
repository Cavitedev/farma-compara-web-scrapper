package farmaciasdirect

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

const firstPage int = 1

var lastPage int = firstPage

func ScrapItems(client *firestore.Client) {

	log.Println(Domain)
	c := colly.NewCollector(
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)

	})

	c.OnHTML("#js-product-list", func(h *colly.HTMLElement) {
		log.Println("Product List")
		if lastPage == firstPage {
			pageItems := h.ChildTexts(".page-item")
			pageStr := pageItems[len(pageItems)-2]
			lastPageI64, err := strconv.ParseInt(pageStr, 10, 32)
			if err != nil {
				log.Println("Could not parse page number")
			} else {
				lastPage = int(lastPageI64)

			}
		}

		h.ForEach(".card-product", func(_ int, e *colly.HTMLElement) {
			item := Item{}
			pageItem := WebsiteItem{}
			pageItem.Url = e.ChildAttr(".card-body>a", "href")
			scrapDetailsPage(&item, &pageItem)
			if item.WebsiteItems == nil {
				item.WebsiteItems = make(map[string]WebsiteItem)
			}
			item.WebsiteItems[websiteName] = pageItem
			firestore_utils.UpdateItem(item, client)
			time.Sleep(50 * time.Millisecond)
		})
	})

	for i := firstPage; i <= lastPage; i++ {
		url := buildPageUrl(i)
		log.Println("Visit Page", i, " url:", url)
		err := c.Visit(url)
		if err != nil {
			log.Printf("Error when visiting %v, err:%v", url, err)
		}
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

func buildPageUrl(pageNum int) string {

	url := fmt.Sprintf("https://%v/corporal-1240?page=%v", Domain, pageNum)
	pageNum++
	return url
}
