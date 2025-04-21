package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"public-service/config"
	"public-service/controller"
	"public-service/repository"
	db "public-service/scripts"
	"public-service/service"
	"public-service/utils"
)

func main() {
	// Load env and initialize common stuff
	config.LoadEnv()
	utils.InitLogger()
	repository.InitDB()
	db.RunMigrations()

	// Initialize repositories
	dbConn := repository.GetDB()
	publicRepo := repository.NewPublicRepository(dbConn)
	appRepo := repository.NewApplicationRepository(dbConn, publicRepo) // Inject PublicRepo

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
		port = "8080" // default fallback
	}
	log.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
