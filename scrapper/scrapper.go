package scrapper

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/scrapper/dosfarma"
	"github.com/cavitedev/go_tuto/scrapper/farmaciaencasa"
	"github.com/cavitedev/go_tuto/scrapper/farmaciasdirect"
	"github.com/cavitedev/go_tuto/scrapper/okfarma"
)

func Scrap(website string, ref *firestore.CollectionRef) {

	log.Println("Hola scrapper")

	switch website {
	case okfarma.Domain:
		okfarma.Scrap(ref)
	case farmaciasdirect.Domain:
		farmaciasdirect.Scrap(ref)
	case dosfarma.Domain:
		dosfarma.Scrap(ref)
	case farmaciaencasa.Domain:
		farmaciaencasa.Scrap(ref)
	}

}
