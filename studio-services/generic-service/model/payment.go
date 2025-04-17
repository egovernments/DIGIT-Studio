package model

type PaymentRequest struct {
	Amount float64 `json:"amount"`
	Payer  string  `json:"payer"`
}
