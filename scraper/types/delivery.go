package types

type Delivery struct {
	Url       string                  `json:"url" firestore:"url"`
	Locations map[string][]PriceRange `json:"locations" firestore:"locations"`
}

type PriceRange struct {
	Price float64 `json:"price" firestore:"price"`
	Min   float64 `json:"min" firestore:"min"`
	Max   float64 `json:"max" firestore:"max"`
}
