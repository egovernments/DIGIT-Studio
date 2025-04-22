package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"public-service/controller"
	"public-service/repository"
	db "public-service/scripts"
	"public-service/service"
	"public-service/utils"
)

func main() {
	utils.InitLogger()

	// Init DB always, migrations optional
	dbConn := repository.InitDB()

	if os.Getenv("FLYWAY_ENABLED") == "true" {
		db.RunMigrations()
	}

	// Initialize repositories
	publicRepo := repository.NewPublicRepository(dbConn)
	appRepo := repository.NewApplicationRepository(dbConn, publicRepo)

	// Initialize services
	appSvc := service.NewApplicationService(appRepo)
	serviceSvc := service.NewPublicService(publicRepo)

	// Initialize controllers
	appCtrl := controller.NewApplicationController(appSvc)
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
