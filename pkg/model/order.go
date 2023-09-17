package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderID string  `json:"order_id"`
	Side    string  `json:"side"`
	Symbol  string  `json:"symbol"`
	Amount  float64 `json:"amount"`
	Price   float64 `json:"price"`
}
