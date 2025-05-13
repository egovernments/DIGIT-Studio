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
	"strconv"
	"strings"
)

type SMSService struct {
	restCallRepo        repository.RestCallRepository
	localizationService *LocalizationService
	kafkaProducer       *producer.PublicServiceProducer
}

func NewSMSService(repo repository.RestCallRepository, localizationService *LocalizationService, kafkaProducer *producer.PublicServiceProducer) *SMSService {
	return &SMSService{
		restCallRepo:        repo,
		localizationService: localizationService,
		kafkaProducer:       kafkaProducer,
	}
}

func (s *SMSService) SendSMS(application model.ApplicationRequest, tenantId string, templateCode string, owners []model.Applicant) (map[string]interface{}, error) {
	localizationMessage := s.localizationService.GetLocalizationMessage(
		application.RequestInfo,
		templateCode,
		tenantId,
	)

	templateMsg := localizationMessage["message"]

	if s.kafkaProducer == nil {
		return nil, fmt.Errorf("Kafka producer is not initialized")
	}

	for _, owner := range owners {
		msg := templateMsg

		if owner.Name != "" {
			msg = strings.ReplaceAll(msg, "{username}", owner.Name)
		}
		if application.Application.ApplicationNumber != "" {
			msg = strings.ReplaceAll(msg, "{applicationNo}", application.Application.ApplicationNumber)
		}
		msg = strings.ReplaceAll(msg, "{taxamout}", "50.00")

		smsRequest := sms.SMSRequest{
			MobileNumber: strconv.FormatInt(owner.MobileNumber, 10),
			Message:      msg,
			Category:     sms.CategoryNotification,
			TenantID:     tenantId,
		}

		smsBytes, err := json.Marshal(smsRequest)
		if err != nil {
			log.Printf("Failed to marshal SMSRequest for owner %v: %v", owner.MobileNumber, err)
			continue
		}

		ctx := context.Background()
		err = s.kafkaProducer.Push(ctx, config.GetEnv("SEND_SMS_TOPIC"), smsBytes)
		if err != nil {
			log.Printf("Failed to push Kafka message for owner %v: %v", owner.MobileNumber, err)
			continue
		}
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Messages sent",
	}, nil
}
