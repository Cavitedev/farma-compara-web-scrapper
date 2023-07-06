package scraper

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/scraper/dosfarma"
	"github.com/cavitedev/go_tuto/scraper/farmaciaencasa"
	"github.com/cavitedev/go_tuto/scraper/farmaciasdirect"
	"github.com/cavitedev/go_tuto/scraper/okfarma"
)

func Scrap(website string, client *firestore.Client, scrapItems bool, scrapDelivery bool) {

	log.Println("Hola scraper")

	switch website {
	case okfarma.Domain:
		okfarma.Scrap(client, scrapItems, scrapDelivery)
	case farmaciasdirect.Domain:
		farmaciasdirect.Scrap(client, scrapItems, scrapDelivery)
	case dosfarma.Domain:
		dosfarma.Scrap(client, scrapItems, scrapDelivery)
	case farmaciaencasa.Domain:
		farmaciaencasa.Scrap(client, scrapItems, scrapDelivery)
	}

}
