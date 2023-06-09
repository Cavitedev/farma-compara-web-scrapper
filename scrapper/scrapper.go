package scrapper

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/scrapper/dosfarma"
	"github.com/cavitedev/go_tuto/scrapper/farmaciaencasa"
	"github.com/cavitedev/go_tuto/scrapper/farmaciasdirect"
	"github.com/cavitedev/go_tuto/scrapper/okfarma"
)

func Scrap(website string, client *firestore.Client, scrapItems bool, scrapDelivery bool) {

	log.Println("Hola scrapper")

	switch website {
	case okfarma.Domain:
		okfarma.Scrap(client, scrapItems, scrapDelivery)
	case farmaciasdirect.Domain:
		farmaciasdirect.Scrap(client)
	case dosfarma.Domain:
		dosfarma.Scrap(client, scrapItems, scrapDelivery)
	case farmaciaencasa.Domain:
		farmaciaencasa.Scrap(client, scrapItems, scrapDelivery)
	}

}
