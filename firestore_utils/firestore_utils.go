package firestore_utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/cavitedev/go_tuto/firestore_utils/transform"
	. "github.com/cavitedev/go_tuto/scrapper/types"
)

func UpdateItem(item Item, col *firestore.CollectionRef) {
	ctx := context.Background()
	id := item.Ref
	doc := col.Doc(id)

	m := transform.ToFirestoreMap(item)

	_, err := doc.Set(ctx, m, firestore.MergeAll)

	if err != nil {
		log.Panicf("Could not insert %v\n", item)
	}
}
