package scrapper

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/scrapper/farmaciasdirect"
	"github.com/cavitedev/go_tuto/scrapper/okfarma"
)

func Scrap(website string, ref *firestore.CollectionRef) {

	fmt.Println("Hola scrapper")

	if website == okfarma.Domain {
		okfarma.Scrap(ref)
	} else if website == "www.farmaciasdirect.com" {
		farmaciasdirect.Scrap()
	}

}
