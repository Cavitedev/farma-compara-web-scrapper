package farmaciasdirect

import (
	"log"

	"cloud.google.com/go/firestore"
)

const Domain string = "www.farmaciasdirect.com"

func Scrap(ref *firestore.CollectionRef) {

	log.Println(Domain)

}
