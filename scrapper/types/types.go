package types

type Item struct {
	Ref       string `json:"ref"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Image     string `json:"image"`
	Url       string `json:"url"`
	Available bool   `json:"available"`
}
