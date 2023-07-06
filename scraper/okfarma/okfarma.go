package okfarma

import (
	"log"

	"cloud.google.com/go/firestore"
)

const websiteName string = "okfarma"
const Domain string = "okfarma.es"

func Scrap(client *firestore.Client, scrapItems bool, scrapDelivery bool) {

	log.Println(Domain)

	if scrapItems {
		ScrapItems(client)
	}

	if scrapDelivery {
		ScrapDelivery(client)
	}

}
