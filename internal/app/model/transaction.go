package model

type Transaction struct {
	Model
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
