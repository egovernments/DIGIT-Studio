package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"public-service/config"
	producer "public-service/kafka/producer"
	"public-service/model"
	_ "public-service/model/demand"
	"public-service/model/sms"
	"public-service/repository"
	"strconv"
	"strings"
)

type SMSService struct {
	restCallRepo        repository.RestCallRepository
	localizationService *LocalizationService
	kafkaProducer       *producer.PublicServiceProducer
	demandService       *DemandService
}

func NewSMSService(repo repository.RestCallRepository, localizationService *LocalizationService, kafkaProducer *producer.PublicServiceProducer, demandService *DemandService) *SMSService {
	return &SMSService{
		restCallRepo:        repo,
		localizationService: localizationService,
		kafkaProducer:       kafkaProducer,
		demandService:       demandService,
	}
}

func (s *SMSService) SendSMS(application model.ApplicationRequest, tenantId string, templateCode string, owners []model.Applicant) (map[string]interface{}, error) {
	localizationMessage := s.localizationService.GetLocalizationMessage(
		application.RequestInfo,
		templateCode,
		tenantId,
	)

	templateMsg := localizationMessage["message"]
	if templateMsg == "" {
		log.Println("Localization message not found for template:", templateCode)
		return nil, fmt.Errorf("template message not found")
	}

	if s.kafkaProducer == nil {
		return nil, fmt.Errorf("Kafka producer is not initialized")
	}

	// Fetch bills and calculate total amount
	bills, err := s.demandService.fetchBill(application)
	if err != nil {
		log.Printf("Error fetching bill for application %s: %v", application.Application.ApplicationNumber, err)
		return nil, err
	}

	var totalAmount float64
	for _, bill := range bills {
		totalAmount += bill.TotalAmount
	}
	amountStr := strconv.FormatFloat(totalAmount, 'f', 2, 64)

	// Loop over all owners to send SMS
	for _, owner := range owners {
		msg := templateMsg

		if owner.Name != "" {
			msg = strings.ReplaceAll(msg, "{PublicService.applicants[0].name}", owner.Name)
		}
		if application.Application.ApplicationNumber != "" {
			msg = strings.ReplaceAll(msg, "{PublicService.applicationNo}", application.Application.ApplicationNumber)
		}
		msg = strings.ReplaceAll(msg, "{Bill.totalAmount}", amountStr)

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

		//err = s.kafkaProducer.Push(ctx, config.GetEnv("SEND_NOTIFICATION_TOPIC"), smsBytes)
		//if err != nil {
		//	log.Printf("Failed to push Kafka message for owner %v: %v", owner.MobileNumber, err)
		//	continue
		//}
	}

	return map[string]interface{}{
		"status":  "success",
		"message": "Messages sent",
	}, nil
}
