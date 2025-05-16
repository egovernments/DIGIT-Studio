package controller

import (
	"Application-Service/model"
	"Application-Service/service"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ApplicationController struct {
	service *service.ApplicationService
}

func NewApplicationController(service *service.ApplicationService) *ApplicationController {
	return &ApplicationController{service: service}
}

func (c *ApplicationController) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var req model.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	res, err := c.service.CreateApplication(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *ApplicationController) SearchApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var criteria model.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&criteria)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Println("inside search", criteria.SearchCriteria)
	ctx := context.Background()
	res, err := c.service.SearchApplication(ctx, criteria.SearchCriteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *ApplicationController) UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var req model.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Println("inside search", req)
	ctx := context.Background()
	res, err := c.service.UpdateApplication(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
