package dosfarma

import (
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	"github.com/cavitedev/go_tuto/scraper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

var deliveryUrl = "https://www.dosfarma.com/envios"

// Solo saca los datos de península cuando pesa menos de 20 KG
func ScrapDelivery(client *firestore.Client) {

	log.Printf("Scrappeando envíos de %v\n", Domain)

	delivery := types.Delivery{}
	delivery.Url = deliveryUrl
	delivery.Locations = make(map[string][]types.PriceRange)

	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)

	var executed bool = false

	c.OnHTML("table.MsoTableGrid>tbody", func(h *colly.HTMLElement) {

		if executed {
			return
		}

		var pricesRange []types.PriceRange = []types.PriceRange{}

		h.ForEach("tr", func(i int, tr *colly.HTMLElement) {

			if i == 1 {

				tr.ForEach("td", func(j int, td *colly.HTMLElement) {
					if j == 3 {
						price := PriceFromTableCell(td)
						pricesRange = append(pricesRange, types.PriceRange{Price: price, Min: 0, Max: 49})
					}
				})
			} else if i == 2 {
				tr.ForEach("td", func(j int, td *colly.HTMLElement) {
					if j == 1 {
						price := PriceFromTableCell(td)
						pricesRange = append(pricesRange, types.PriceRange{Price: price, Min: 49.01})
					}
				})
			}

		})
		delivery.Locations["spain"] = pricesRange
		delivery.Locations["portugal"] = pricesRange
		executed = true
	})

	c.Visit(deliveryUrl)

	firestore_utils.UpdateDelivery(delivery, client, websiteName)

}

func PriceFromTableCell(td *colly.HTMLElement) float64 {
	text := td.Text

	var price float64

	if utils.IsNumber(text) {
		price = utils.ParseSpanishNumberStrToNumber(text)
	}
	return price
}

func InnerRowDelivery(j int, td *colly.HTMLElement, pricesRange []types.PriceRange, keys []string) types.PriceRange {

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
func HeaderRowDelivery(j int, td *colly.HTMLElement, pricesRange []types.PriceRange) []types.PriceRange {
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

		pricesRange = append(pricesRange, types.PriceRange{
			Min: minPrice,
			Max: maxPrice,
		})
	}
	return pricesRange
}
