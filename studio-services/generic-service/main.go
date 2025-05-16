package main

import (
	db "egov-generic-service/scripts"
	"log"
	"net/http"

	"egov-generic-service/controller"
	"egov-generic-service/repository"
	"egov-generic-service/utils"
	"github.com/gorilla/mux"
)

func main() {
	utils.InitLogger()
	repository.InitDB()
	db.RunMigrations()
	router := mux.NewRouter()

	router.HandleFunc("/generic-service/v1/service/_init", controller.InitServiceHandler).Methods("POST")
	router.HandleFunc("/generic-service/v1/service/_update", controller.UpdateServiceHandler).Methods("POST")
	router.HandleFunc("/generic-service/v1/service/_search", controller.SearchServiceHandler).Methods("POST")

	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
