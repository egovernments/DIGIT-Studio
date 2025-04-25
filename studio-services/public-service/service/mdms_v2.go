package service

import (
	"fmt"
	"time"

	"github.com/Priyansuvaish/digit_client/models"
	"github.com/Priyansuvaish/digit_client/services"
	"github.com/google/uuid"
)

// CreateMDMS creates an MDMS entry with the given auth token
func CreateMDMS(token string) error {
	schemaCode := "abc"

	// Initialize RequestConfig properly
	requestConfig := (&models.RequestConfig{}).GetInstance()
	requestConfig.Initialize(
		"digit",                  // apiID
		"1.0",                    // version
		token,                    // authToken
		map[string]interface{}{}, // userInfo
		"device123",              // did
		"key123",                 // key
		uuid.New().String(),      // msgID
		"requester123",           // requesterID
		uuid.New().String(),      // correlationID
		"authorize",              // action
		time.Now().UnixMilli(),   // ts
	)

	// Build MDMS data
	MDMS := models.MdmsBuilder().
		WithTenantID("LMN").
		WithSchemeCode("Owner.cardetails").
		WithData(map[string]interface{}{
			"ownerName":     "mustakim N kouji",
			"contactNumber": []int{123, 9090},
			"address":       "",
			"state":         "",
			"RTO":           105,
			"car": map[string]interface{}{
				"moduleName": "122",
				"color":      "blue",
				"carNumber":  123,
				"engine": map[string]string{
					"Ename":    "",
					"capacity": "",
					"power":    "",
					"average":  "123",
				},
			},
		}).
		WithIsActive(true).
		Build()

	// Create MDMS service instance
	mdmsService := services.NewMdmsV2Service(nil)

	// Call CreateMDMS method
	createdMDMS, err := mdmsService.CreateMDMS(schemaCode, MDMS, nil)
	if err != nil {
		return err
	}

	// Optionally do something with createdMDMS
	fmt.Println("MDMS created:", createdMDMS)

	return nil
}
