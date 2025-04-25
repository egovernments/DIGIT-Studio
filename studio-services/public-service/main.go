package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"public-service/controller"
	"public-service/repository"
	db "public-service/scripts"
	"public-service/service"
	"public-service/utils"
	"strings"

	"github.com/Priyansuvaish/digit_client/config"
	"github.com/Priyansuvaish/digit_client/digit_init"

	"github.com/gorilla/mux"
)

func main() {

	//to generate the Auth token
	GenerateAndSaveAuthToken()

	// configre the url and auth token globally
	config.GetGlobalConfig().Initialize(
		"https://sandbox.digit.org",
		os.Getenv("AUTH_TOKEN"),
	)

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
	appCtrl := controller.NewApplicationController(appSvc, service.NewWorkflowIntegrator())
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

func GenerateAndSaveAuthToken() {
	// Step 1: Authenticate and get token
	token, err := digit_init.Authenticate()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}
	fmt.Println("Received token:", token)

	// Step 2: Update or append AUTH_TOKEN in .env file
	envFile := ".env"

	// Read existing lines
	lines, err := readLines(envFile)
	if err != nil {
		// If file doesn't exist, create new
		if os.IsNotExist(err) {
			lines = []string{}
		} else {
			log.Fatalf("Failed to read .env file: %v", err)
		}
	}

	// Update or add AUTH_TOKEN line
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "AUTH_TOKEN=") {
			lines[i] = fmt.Sprintf("AUTH_TOKEN=%s", token)
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, fmt.Sprintf("AUTH_TOKEN=%s", token))
	}

	// Write back to .env file
	err = writeLines(lines, envFile)
	if err != nil {
		log.Fatalf("Failed to write to .env file: %v", err)
	}

	fmt.Println(".env file updated successfully")
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return w.Flush()
}
