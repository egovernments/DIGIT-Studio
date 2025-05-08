package model

type BillingRequest struct {
	BillId string  `json:"billId"`
	Amount float64 `json:"amount"`
}
