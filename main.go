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

var domain string = "okfarma.es"

func main() {

	var err error

	ctx = context.Background()
	sa := option.WithCredentialsFile("secrets/functions-sa.json")
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

	// _, _, err = ref.Add(ctx, map[string]interface{}{
	// 	"first": "Ada",
	// 	"last":  "Lovelace",
	// 	"born":  1815,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed adding alovelace: %v", err)
	// }

}
