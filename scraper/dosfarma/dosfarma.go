package dosfarma

import (
	"cloud.google.com/go/firestore"
	"log"
)

const websiteName string = "dosfarma"
const Domain string = "www.dosfarma.com"

func Scrap(client *firestore.Client, scrapItems bool, scrapDelivery bool) {

	log.Println(Domain)

	if scrapItems {
		ScrapItems(client)
	}

	if scrapDelivery {
		ScrapDelivery(client)
	}

}
