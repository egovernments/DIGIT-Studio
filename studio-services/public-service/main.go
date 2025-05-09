package main

import (
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"public-service/config"
	"public-service/controller"
	producer "public-service/kafka"
	"public-service/repository"
	db "public-service/scripts"
	"public-service/service"
	"public-service/utils"
)

func main() {
	utils.InitLogger()

	// Init DB always, migrations optional
	dbConn := repository.InitDB()
	config.LoadEnv()

	if os.Getenv("FLYWAY_ENABLED") == "true" {
		db.RunMigrations()
	}
	writerFunc := func(topic string) *kafka.Writer {
		return kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{config.GetEnv("KAFKA_BOOTSTRAP_SERVERS")}, // Example: "localhost:9092"
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		})
	}
	kafkaProducer := producer.NewPublicServiceProducer(writerFunc)
	// Initialize repositories
	publicRepo := repository.NewPublicRepository(dbConn)
	appRepo := repository.NewApplicationRepository(dbConn, publicRepo, kafkaProducer)
	restRepo := repository.NewRestCallRepository()
	// Initialize services
	demandSvc := service.NewDemandService(restRepo)
	individualSvc := service.NewIndividualService(restRepo)
	mdmsSvc := service.NewMDMSService(restRepo)
	mdmsv2sSvc := service.NewMDMSV2Service(restRepo)
	idgenSvc := service.NewIdGenService(restRepo)
	enrichSvc := service.NewEnrichmentService(individualSvc, demandSvc, mdmsSvc, mdmsv2sSvc, idgenSvc)
	appSvc := service.NewApplicationService(appRepo, enrichSvc)
	serviceSvc := service.NewPublicService(publicRepo)
	// Initialize controllers
	appCtrl := controller.NewApplicationController(appSvc, service.NewWorkflowIntegrator(mdmsv2sSvc), individualSvc, enrichSvc)
	serviceCtrl := controller.NewServiceController(serviceSvc)

	// Setup router
	router := mux.NewRouter()

	// Service routes
	router.HandleFunc("/public-service/v1/service", serviceCtrl.CreateServiceHandler).Methods("POST")
	router.HandleFunc("/public-service/v1/service", serviceCtrl.SearchServiceHandler).Methods("GET")
	router.HandleFunc("/public-service/v1/service/{serviceCode}", serviceCtrl.UpdateServiceHandler).Methods("PUT")

	// Application routes
	router.HandleFunc("/public-service/v1/application/{serviceCode}", appCtrl.CreateApplicationHandler).Methods("POST")
	router.HandleFunc("/public-service/v1/application/{serviceCode}", appCtrl.SearchApplicationHandler).Methods("GET")
	router.HandleFunc("/public-service/v1/application/{serviceCode}/{applicationId}", appCtrl.UpdateApplicationHandler).Methods("PUT")

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf(" Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
