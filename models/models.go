package models

type Stock struct {
	ID    string  `json:"stockid"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Company string `josn:"company"`
}
