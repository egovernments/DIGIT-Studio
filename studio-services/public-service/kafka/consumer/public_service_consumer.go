package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"public-service/config"
	"public-service/model"
	"public-service/model/payment"
	"public-service/service"

	"github.com/segmentio/kafka-go"
)

func ConsumePayments(applicationService *service.ApplicationService) {
	topic := os.Getenv("KAFKA_TOPICS_PAYMENT_CREATE_NAME")
	if topic == "" {
		log.Fatal("❌ KAFKA_TOPICS_PAYMENT_CREATE_NAME is not set")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.GetEnv("KAFKA_BOOTSTRAP_SERVERS")},
		GroupID:  "public-service-group",
		Topic:    topic,
		MaxBytes: 10e6, // 10MB max per message
	})
	defer r.Close()

	log.Printf("📡 Kafka consumer started on topic: %s", topic)

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		var paymentReq payment.PaymentRequest
		if err := json.Unmarshal(m.Value, &paymentReq); err != nil {
			log.Printf("❌ Failed to unmarshal payment request: %v", err)
			continue
		}

		if len(paymentReq.Payment.PaymentDetails) == 0 {
			log.Printf("⚠️ No payment details found in message: %+v", paymentReq)
			continue
		}

		criteria := model.SearchCriteria{
			TenantId:         paymentReq.Payment.TenantID,
			ApplicationNumber: paymentReq.Payment.PaymentDetails[0].Bill.ConsumerCode,
			BusinessService:  paymentReq.Payment.PaymentDetails[0].BusinessService,
		}

		searchRes, err := applicationService.SearchApplication(context.Background(), criteria)
		if err != nil || len(searchRes.Application) == 0 {
			log.Printf("❌ Application not found for payment: %+v, error: %v", criteria, err)
			continue
		}

		application := searchRes.Application[0]
		application.Workflow.Action = "PAY"

		appReq := model.ApplicationRequest{
			RequestInfo: paymentReq.RequestInfo,
			Application: application,
		}

		log.Printf("📩 Payment received for application [%s] on topic [%s]", application.ApplicationNumber, m.Topic)

		_, err = applicationService.UpdateApplication(context.Background(), appReq, application.ServiceCode, application.Id.String())
		if err != nil {
			log.Printf("❌ Failed to update application after payment: %v", err)
			continue
		}

		log.Printf("✅ Application [%s] updated successfully after payment", application.ApplicationNumber)
	}
}
