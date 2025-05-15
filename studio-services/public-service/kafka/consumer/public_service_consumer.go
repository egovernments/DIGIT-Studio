package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"public-service/config"
	"public-service/model/payment"
	"public-service/service"
    "public-service/model"
	"github.com/segmentio/kafka-go"
)

func ConsumePayments(workflowIntegrator *service.WorkflowIntegrator) {
	topic := os.Getenv("KAFKA_TOPICS_PAYMENT_CREATE_NAME")
	if topic == "" {
		log.Fatal("KAFKA_TOPICS_PAYMENT_CREATE_NAME is not set")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.GetEnv("KAFKA_BOOTSTRAP_SERVERS")},
		GroupID:  "public-service-group",
		Topic:    topic,
		MaxBytes: 10e6, // 10MB max per message
	})
	defer r.Close()

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
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

		applicationRequest := &model.ApplicationRequest{
			RequestInfo: paymentReq.RequestInfo,
			Application: model.Application{
				ApplicationNumber: paymentReq.Payment.PaymentDetails[0].Bill.ConsumerCode,
				TenantId:          paymentReq.Payment.TenantID,
				Workflow:          model.WorkFlow{Action: "PAY"},
			},
		}

		log.Printf("📩 Received payment on topic [%s] for application [%s]",
			m.Topic, applicationRequest.Application.ApplicationNumber,
		)

		if err := workflowIntegrator.CallWorkflow(applicationRequest); err != nil {
			log.Printf("❌ Workflow call failed: %v", err)
			continue
		}

		if err := r.CommitMessages(context.Background(), m); err != nil {
			log.Printf("❌ Failed to commit message: %v", err)
		}
	}
}


