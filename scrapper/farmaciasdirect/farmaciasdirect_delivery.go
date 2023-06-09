package farmaciasdirect

import (
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	"github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

var deliveryUrl = "https://www.farmaciasdirect.com/informacion/envios-12"

func ScrapDelivery(client *firestore.Client) {

	log.Printf("Scrappeando envíos de %v\n", Domain)

	delivery := types.Delivery{}
	delivery.Url = deliveryUrl
	delivery.Locations = make(map[string][]types.PriceRange)

	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)

	var executed bool = false

	c.OnHTML("#content>table>tbody", func(h *colly.HTMLElement) {

		if executed {
			return
		}
		var keyLists []string = []string{"spain", "portugal", "balearic"}

		//Rango de precios sin valor

		var pricesRange []types.PriceRange = []types.PriceRange{}

		h.ForEach("tr", func(i int, tr *colly.HTMLElement) {

			//Conjutno de lugares en cada fila
			if i == 1 {

				tr.ForEach("td", func(j int, td *colly.HTMLElement) {

					if j > 0 {
						text := td.Text
						var addPriceRange types.PriceRange
						var price float64

						if utils.IsNumber(text) {
							price = utils.ParseSpanishNumberStrToNumber(text)
						}

						if strings.Contains(text, "superiores a") {
							addPriceRange = types.PriceRange{Min: price}
						} else if strings.Contains(text, "inferiores a") {
							addPriceRange = types.PriceRange{Max: price}
						} else {
							log.Fatalf("No se pudo parsear el texto: %v\n", text)
						}

						pricesRange = append(pricesRange, addPriceRange)
					}
				})
			} else if i == 2 {
				{

					tr.ForEach("td", func(j int, td *colly.HTMLElement) {

						if j > 0 {
							text := td.Text
							var price float64

							if utils.IsNumber(text) {
								price = utils.ParseSpanishNumberStrToNumber(text)
							}

							pricesRange[j-1].Price = price
						}
					})
				}
			}

		})

		for _, key := range keyLists {
			delivery.Locations[key] = pricesRange
		}

		executed = true
	})

	c.Visit(deliveryUrl)

	firestore_utils.UpdateDelivery(delivery, client, websiteName)

}
