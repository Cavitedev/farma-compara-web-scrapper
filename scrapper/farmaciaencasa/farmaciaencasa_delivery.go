package farmaciaencasa

import (
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	. "github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

var deliveryUrl = "https://www.farmaciaencasaonline.es/farmacia-en-casa-online-tarifas-y-condiciones-de-envio"

func ScrapDelivery(client *firestore.Client) {

	log.Printf("Scrappeando envíos de %v\n", Domain)

	delivery := Delivery{}
	delivery.Url = deliveryUrl
	delivery.Locations = make(map[string][]PriceRange)

	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)

	var executed bool = false

	c.OnHTML(".main>div>table>tbody", func(h *colly.HTMLElement) {

		if executed {
			return
		}
		var keyLists [][]string = [][]string{}

		//Rango de precios sin valor
		var tmpPricesRange []PriceRange
		var pricesRange []PriceRange

		h.ForEach("tr", func(i int, tr *colly.HTMLElement) {

			//Conjutno de lugares en cada fila
			if i == 0 || i == 4 {
				tmpPricesRange = []PriceRange{}
				tr.ForEach("th", func(j int, td *colly.HTMLElement) {
					tmpPricesRange = HeaderRowDelivery(j, td, tmpPricesRange)
				})
			} else {
				var keys []string
				tr.ForEach("td", func(j int, td *colly.HTMLElement) {

					if j == 0 {
						keys = getRegions(td)
						keyLists = append(keyLists, keys)
					} else {
						addPriceRange := InnerRowDelivery(j, td, tmpPricesRange, keys)

						pricesRange = append(pricesRange, addPriceRange)

					}
				})
			}

		})

		for i, keyList := range keyLists {
			for _, key := range keyList {
				delivery.Locations[key] = pricesRange[i*3 : i*3+3]
			}
		}

		executed = true
	})

	c.Visit(deliveryUrl)

	firestore_utils.UpdateDelivery(delivery, client, websiteName)

}

func InnerRowDelivery(j int, td *colly.HTMLElement, pricesRange []PriceRange, keys []string) PriceRange {

	text := td.Text

	var price float64

	if utils.IsNumber(text) {
		price = utils.ParseSpanishNumberStrToNumber(text)
	}

	pricesRange[j-1].Price = price
	return pricesRange[j-1]

}

func getRegions(td *colly.HTMLElement) []string {

	var keys []string = []string{}

	region := td.Text

	switch region {
	case "Madrid (Comunidad)":
		keys = append(keys, "madrid")
	case "España (Península)":
		keys = append(keys, "spain")
	case "Portugal (Península)":
		keys = append(keys, "portugal")
	case "Baleares":
		keys = append(keys, "balearic")
	case "Formentera(1)":
		keys = append(keys, "formentera")
	case "Canarias, Ceuta y Melilla(2)":
		keys = append(keys, "canary")
		keys = append(keys, "melilla")
		keys = append(keys, "ceuta")
	default:
		log.Fatalf("No se entiende el lugar de envio: %v", region)
	}
	return keys
}
func HeaderRowDelivery(j int, td *colly.HTMLElement, pricesRange []PriceRange) []PriceRange {
	if j > 0 {
		text := td.Text

		splittedText := strings.SplitN(text, "–", 2)

		var minPrice float64
		var maxPrice float64

		if utils.IsNumber(splittedText[0]) {
			minPrice = utils.ParseSpanishNumberStrToNumber(splittedText[0])
		}
		if len(splittedText) >= 2 && utils.IsNumber(splittedText[1]) {
			maxPrice = utils.ParseSpanishNumberStrToNumber(splittedText[1]) - 0.01
		} else {
			minPrice = utils.ParseSpanishNumberStrToNumber(splittedText[0])
		}

		pricesRange = append(pricesRange, PriceRange{
			Min: minPrice,
			Max: maxPrice,
		})
	}
	return pricesRange
}
