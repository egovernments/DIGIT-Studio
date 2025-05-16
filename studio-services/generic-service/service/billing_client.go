package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"egov-generic-service/model"
)

func CallBillingService(billingRequest model.BillingRequest) error {
	url := "http://billing-service/billing/_bill"
	body, _ := json.Marshal(billingRequest)

	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	return err
}
