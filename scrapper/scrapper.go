package scrapper

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/scrapper/farmaciasdirect"
	"github.com/cavitedev/go_tuto/scrapper/okfarma"
)

func Scrap(website string, ref *firestore.CollectionRef) {

	log.Println("Hola scrapper")

	if website == okfarma.Domain {
		okfarma.Scrap(ref)
	} else if website == farmaciasdirect.Domain {
		farmaciasdirect.Scrap(ref)
	}

}
