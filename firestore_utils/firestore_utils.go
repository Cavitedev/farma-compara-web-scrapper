package firestore_utils

import (
	"context"
	"log"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils/transform"
	"github.com/cavitedev/go_tuto/scraper/types"
)

func UpdateItem(item types.Item, client *firestore.Client) {

	col := client.Collection("items")

	ctx := context.Background()
	id := strings.Replace(item.Ref, "/", "_", -1)

	m1 := regexp.MustCompile(`Ref\.|p-`)
	id = m1.ReplaceAllString(id, "")

	m2 := regexp.MustCompile(`Ref\.|(\d*)\..*`)

	id = m2.ReplaceAllString(id, "$1")

	item.Ref = id
	doc := col.Doc(id)

	m := transform.ToFirestoreMap(item)

	_, err := doc.Set(ctx, m, firestore.MergeAll)

	if err != nil {
		log.Printf("Error: %v Could not insert %v\n", err, item)
	}
}

func UpdateDelivery(delivery types.Delivery, client *firestore.Client, docKey string) {
	ctx := context.Background()
	col := client.Collection("delivery_fees")

	m := transform.ToFirestoreMap(delivery)
	doc := col.Doc(docKey)

	_, err := doc.Set(ctx, m, firestore.MergeAll)

	if err != nil {
		log.Printf("Error: %v Could not insert %v\n", err, delivery)
	}
}
