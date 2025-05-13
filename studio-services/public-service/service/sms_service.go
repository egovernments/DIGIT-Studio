package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"public-service/config"
	producer "public-service/kafka"
	"public-service/model"
	"public-service/model/sms"
	"public-service/repository"
)

type SMSService struct {
	restCallRepo        repository.RestCallRepository
	localizationService LocalizationService
	kafkaProducer       *producer.PublicServiceProducer
}

func NewSMSService(repo repository.RestCallRepository, localizationService LocalizationService, kafkaProducer *producer.PublicServiceProducer) *SMSService {
	return &SMSService{
		restCallRepo:        repo,
		localizationService: localizationService,
		kafkaProducer:       kafkaProducer,
	}
}

func (s *SMSService) SendSMS(requestInfo model.RequestInfo, tenantId string, templateCode string, owner model.User) (map[string]interface{}, error) {
	localizationMessage := s.localizationService.GetLocalizationMessage(
		requestInfo,
		templateCode,
		tenantId,
	)

	message := localizationMessage["message"]

	// Construct SMSRequest
	smsRequest := sms.SMSRequest{
		MobileNumber: owner.MobileNumber,
		Message:      message,
		Category:     sms.CategoryTransaction, // Ensure enum is defined like: const CategoryTransaction = "TRANSACTION"
		TenantID:     tenantId,
		// Optionally: ExpiryTime: time.Now().Add(...).Unix(),
	}

	smsBytes, err := json.Marshal(smsRequest)
	if err != nil {
		log.Printf("Failed to marshal SMSRequest: %v", err)
		return nil, err
	}

	// Push to Kafka if producer is available
	if s.kafkaProducer != nil {
		ctx := context.Background()
		err := s.kafkaProducer.Push(ctx, config.GetEnv("SEND_SMS_TOPIC"), smsBytes)
		if err != nil {
			log.Printf("Failed to push Kafka message: %v", err)
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Kafka producer is not initialized")
	}

	return map[string]interface{}{
		"status":  "success",
		"message": message,
	}, nil
}
