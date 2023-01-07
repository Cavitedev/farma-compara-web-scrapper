package types

import "time"

type Item struct {
	Ref          string                 `json:"ref" firestore:"ref"`
	WebsiteItems map[string]WebsiteItem `json:"website_items" firestore:"website_items"`
}

type WebsiteItem struct {
	Name       string    `json:"name" firestore:"name,omitempty"`
	Price      float64   `json:"price" firestore:"price"`
	Image      string    `json:"image" firestore:"image,omitempty"`
	Url        string    `json:"url" firestore:"url"`
	Available  bool      `json:"available" firestore:"available"`
	LastUpdate time.Time `json:"last_update" firestore:"last_update"`
}
