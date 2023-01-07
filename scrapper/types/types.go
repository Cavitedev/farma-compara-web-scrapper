package types

import "time"

type Item struct {
	Ref      string     `json:"ref"`
	PageItem []PageItem `json:"page_item"`
}

type PageItem struct {
	Website    string    `json:"website"`
	Name       string    `json:"name"`
	Price      float32   `json:"price"`
	Image      string    `json:"image"`
	Url        string    `json:"url"`
	Available  bool      `json:"available"`
	LastUpdate time.Time `json:"last_update"`
}
