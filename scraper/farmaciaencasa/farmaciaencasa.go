package farmaciaencasa

import (
	"log"

	"cloud.google.com/go/firestore"
)

const websiteName string = "farmaciaencasa"
const Domain string = "www.farmaciaencasaonline.es"

func Scrap(client *firestore.Client, scrapItems bool, scrapDelivery bool) {

	log.Println(Domain)

	if scrapItems {
		ScrapItems(client)
	}

	if scrapDelivery {
		ScrapDelivery(client)
	}

}
