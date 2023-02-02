package firestore_utils

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils/transform"
	"github.com/cavitedev/go_tuto/scrapper/types"
)

func UpdateItem(item types.Item, col *firestore.CollectionRef) {
	ctx := context.Background()
	id := strings.Replace(item.Ref, "/", "_", -1)
	doc := col.Doc(id)

	m := transform.ToFirestoreMap(item)

	_, err := doc.Set(ctx, m, firestore.MergeAll)

	if err != nil {
		log.Printf("Error: %v Could not insert %v\n", err, item)
	}
}
