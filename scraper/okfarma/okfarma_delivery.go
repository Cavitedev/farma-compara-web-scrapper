package okfarma

import (
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	. "github.com/cavitedev/go_tuto/scraper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

var deliveryUrl = "https://okfarma.es/envio"

func ScrapDelivery(client *firestore.Client) {

	log.Printf("Scrappeando envíos de %v\n", Domain)

	delivery := Delivery{}
	delivery.Url = deliveryUrl
	delivery.Locations = make(map[string][]PriceRange)

	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)

	c.OnHTML(".table-bordered", func(h *colly.HTMLElement) {
		var key string
		var pricesRange []PriceRange = []PriceRange{}
		h.ForEach("tr", func(i int, tr *colly.HTMLElement) {

			if i == 0 {
				key = HeaderRowDelivery(tr, key)
			} else if i == 1 {
				tr.ForEach("td", func(j int, td *colly.HTMLElement) {
					pricesRange = InnerRowDelivery(j, td, pricesRange)
				})
			} else if i == 2 {
				tr.ForEach("td", func(j int, td *colly.HTMLElement) {
					if j > 0 && j < 5 {
						text := td.Text
						var price float64
						if utils.IsNumber(text) {
							price = utils.ParseSpanishNumberStrToNumber(text)
						}
						pricesRange[j-1].Price = price
					}
				})
			}
			if key != "" {
				delivery.Locations[key] = pricesRange

			}

		})

		h.Attr("class")

	})

	c.Visit(deliveryUrl)

	firestore_utils.UpdateDelivery(delivery, client, websiteName)

}

func InnerRowDelivery(j int, td *colly.HTMLElement, pricesRange []PriceRange) []PriceRange {
	if j > 0 && j < 5 {
		text := td.ChildText("strong")

		splittedText := strings.SplitN(text, "€", 2)

		var minPrice float64
		var maxPrice float64

		if utils.IsNumber(splittedText[0]) {
			minPrice = utils.ParseSpanishNumberStrToNumber(splittedText[0])
		}
		if utils.IsNumber(splittedText[1]) {
			maxPrice = utils.ParseSpanishNumberStrToNumber(splittedText[1])
		}

		pricesRange = append(pricesRange, PriceRange{
			Min: minPrice,
			Max: maxPrice,
		})
	}
	return pricesRange
}

func HeaderRowDelivery(tr *colly.HTMLElement, key string) string {
	header := tr.ChildText(".heading-box>td>strong")

	switch header {
	case "ESPAÑA PENINSULAR":
		key = "spain"
	case "BALEARES":
		key = "balearic"
	case "ISLAS CANARIAS, CEUTA Y MELILLA":
	default:
		log.Fatalf("No se entiende el lugar de envio: %v", header)
	}
	return key
}
