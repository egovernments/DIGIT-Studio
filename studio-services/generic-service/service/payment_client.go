package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"egov-generic-service/model"
)

func CallPaymentService(paymentRequest model.PaymentRequest) error {
	url := "http://payment-service/payment/_pay"
	body, _ := json.Marshal(paymentRequest)

	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	return err
}
