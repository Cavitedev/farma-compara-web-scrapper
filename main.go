package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	scraper "github.com/cavitedev/go_tuto/scraper"
	"google.golang.org/api/option"
)

var client *firestore.Client
var ctx context.Context

// var domain string = "www.farmaciasdirect.com"

// var domain string = "okfarma.es"

var domain string = "www.dosfarma.com"

// var domain string = "www.farmaciaencasaonline.es"

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

	fmt.Println("GO")

	scraper.Scrap(domain, client, scrapItems, scrapDelivery)

}
