/*
=== Documentación de la línea de comandos ===

Este programa acepta los siguientes argumentos de la línea de comandos:

- Dominio: Dominio de la página web. Elegir entre: okfarma.es, www.farmaciasdirect.com, www.dosfarma.com, www.farmaciaencasaonline.es y all (para todas las páginas web)
- ScrappearProductos: true si se quiere scrappear los productos, falso en cualquier otro caso
- ScrappearPreciosEnvio: true si se quiere scrappear los precios de envío, falso en cualquier otro caso

Por defecto se utilizan las variables preasignadas

Para ejecutar el programa, utiliza la siguiente sintaxis:

    go run main.go  [Dominio] [ScrappearProductos] [ScrappearPreciosEnvio]

Ejemplos de uso:

    go run main.go

    go run main.go all

	go run main.go okfarma.es true false
=========================================
*/

package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	scraper "github.com/cavitedev/go_tuto/scraper"
	"google.golang.org/api/option"
)

var client *firestore.Client
var ctx context.Context

// Comentar y descomentar para elegir el dominio si no se usan parámetros

// var domain string = "www.farmaciasdirect.com"

// var domain string = "okfarma.es"

var domain string = "www.dosfarma.com"

// var domain string = "www.farmaciaencasaonline.es"

// var domain string = "all"

var scrapItems bool = true
var scrapDelivery bool = true

func main() {

	//Arguments
	if len(os.Args) > 1 {
		domain = os.Args[1]
	}
	if len(os.Args) > 2 {
		scrapItems = os.Args[2] == "true"
	}
	if len(os.Args) > 3 {
		scrapDelivery = os.Args[3] == "true"
	}

	ctx = context.Background()
	sa := option.WithCredentialsFile("secrets/local-functions.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	scraper.Scrap(domain, client, scrapItems, scrapDelivery)

}
