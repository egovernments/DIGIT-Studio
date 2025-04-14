package main

import (
	"Application-Service/config"
	"Application-Service/controller"
	"Application-Service/repository"
	db "Application-Service/scripts"
	"Application-Service/service"
	"Application-Service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	config.LoadEnv()
	utils.InitLogger()
	repository.InitDB()
	db.RunMigrations()
	_ = mux.NewRouter()

	r := gin.Default()

	repo := repository.NewApplicationRepository(repository.GetDB())
	svc := service.NewApplicationService(repo)
	ctrl := controller.NewApplicationController(svc)

	http.HandleFunc("/application-service/v1/_apply", ctrl.CreateApplicationHandler)
	http.HandleFunc("/application-service/v1/_search", ctrl.SearchApplicationHandler)
	http.HandleFunc("/application-service/v1/_update", ctrl.UpdateApplicationHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	r.Run(":" + os.Getenv("SERVER_PORT"))
}
