package farmaciasdirect

import (
	"log"

	"cloud.google.com/go/firestore"
)

const websiteName string = "farmaciasdirect"
const Domain string = "www.farmaciasdirect.com"

func Scrap(client *firestore.Client, scrapItems bool, scrapDelivery bool) {

	log.Println(Domain)

	if scrapItems {
		ScrapItems(client)
	}

	if scrapDelivery {
		ScrapDelivery(client)
	}

}
