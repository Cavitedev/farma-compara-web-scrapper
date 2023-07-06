package farmaciaencasa

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils"
	"github.com/cavitedev/go_tuto/scrapper/types"
	"github.com/cavitedev/go_tuto/utils"
	"github.com/gocolly/colly/v2"
)

var lastPage int = 5
var page int = 1

func ScrapItems(client *firestore.Client) {

	items := []types.Item{}

	// Instanciar colly para que navegue por ese dominio
	c := colly.NewCollector(
		// Se puede programar para que asíncronamente pida distintas páginas webs
		// colly.Async(true),
		colly.AllowedDomains(Domain),
	)

	// Loggear para ver el progreso
	c.OnResponse(func(r *colly.Response) {
		log.Printf("Visit URL:%v\n", r.Request.URL)
	})
	// Tiempo máximo de la petición al cargar el HTML para que no falle por lentitud
	c.SetRequestTimeout(30 * time.Second)

	// Ver las páginas que quedan
	c.OnHTML(".pages", func(h *colly.HTMLElement) {
		pagesLi := h.ChildTexts("li>a")

		// El último elemento es el botón de siguiente, el penúltimo es el número de página
		lastPageLi := pagesLi[len(pagesLi)-2]
		// Elimina todo lo que no sea un número
		lastPageLi = utils.NumberRegexString(lastPageLi)
		lastPageI64, err := strconv.ParseInt(lastPageLi, 10, 32)
		if err != nil {
			log.Println("Error parsing " + lastPageLi)
		}
		// Actualizar la última página
		if lastPageI64 > int64(lastPage) {
			lastPage = int(lastPageI64)
		}

	})

	// Ver las productos de la página
	c.OnHTML(".product-items", func(h *colly.HTMLElement) {

		h.ForEach(".product-item", func(_ int, e *colly.HTMLElement) {

			item := types.Item{}
			pageItem := types.WebsiteItem{}
			pageItem.Url = e.ChildAttr(".product", "href")
			// Una vez tomada la URL delega el scrapeo de la página de detalles
			scrapDetailsPage(&item, &pageItem)
			if item.WebsiteItems == nil {
				item.WebsiteItems = make(map[string]types.WebsiteItem)
			}
			item.WebsiteItems[websiteName] = pageItem
			items = append(items, item)
			// Actualiza la base de datos
			firestore_utils.UpdateItem(item, client)
			// Espera una cantidad de tiempo para no sobrecargar el servidor
			time.Sleep(50 * time.Millisecond)

		})

	})

	// iterar por todas las páginas
	for page != lastPage+1 {
		c.Visit(fmt.Sprintf("https://www.farmaciaencasaonline.es/corporal/cuerpo?p=%v", page))
		page++
	}

}

var productsVisited int = 0

// Página de detalles
func scrapDetailsPage(item *types.Item, pageItem *types.WebsiteItem) {
	c := colly.NewCollector(
		colly.AllowedDomains(Domain),
	)
	c.OnResponse(func(r *colly.Response) {
		productsVisited++
		log.Printf("Visit %d URL:%v\n", productsVisited, r.Request.URL)

	})

	// La imágen esta fuera del elemento que contiene lo demás
	c.OnHTML(".gallery-placeholder__image", func(h *colly.HTMLElement) {
		pageItem.Image = h.Attr("src")
	})

	c.OnHTML(".product-info-main", func(h *colly.HTMLElement) {

		// Textos en divs dentro de la clase sku
		references := h.ChildTexts(".sku>div")
		if len(references) == 0 {
			return
		}
		item.Ref = references[0]

		// Subir el tiempo en el que se consultó la página
		currentTime := time.Now()
		pageItem.LastUpdate = currentTime

		pageItem.Name = h.ChildText(".page-title>span")

		price := h.ChildText(".price")
		// Cuando el precio está rebajado el precio se muestra en otro lugar
		if price == "" {
			price = h.ChildText(".product-type-data>div>.regular-price>.price")
		}

		// Parsear el precio desde el español y con el símbolo del euro a double
		pageItem.Price = utils.ParseSpanishNumberStrToNumber(price)

		availableTexts := h.ChildTexts(".stock>span")
		// Si el texto es disponible entonces está disponible
		pageItem.Available = availableTexts[0] == "Disponible"

	})

	c.Visit(pageItem.Url)
}
