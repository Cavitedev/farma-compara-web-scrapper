package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/cavitedev/go_tuto/scrapper"
	"google.golang.org/api/option"
)

var client *firestore.Client
var ctx context.Context

// var domain string = "www.farmaciasdirect.com"

var domain string = "okfarma.es"

// var domain string = "www.dosfarma.com"

// var domain string = "www.farmaciaencasaonline.es"

func main() {

	ctx = context.Background()
	sa := option.WithCredentialsFile("secrets/farma-functions-sa.json")
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
	ref := client.Collection("items")
	scrapper.Scrap(domain, ref)

}
