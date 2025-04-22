package controller

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"public-service/model"
	"public-service/service"
	"strings"
)

type ApplicationController struct {
	service *service.ApplicationService
}

func NewApplicationController(service *service.ApplicationService) *ApplicationController {
	return &ApplicationController{service: service}
}
func (c *ApplicationController) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	serviceCode := mux.Vars(r)["serviceCode"]

	// Check if serviceCode is missing
	if serviceCode == "" {
		http.Error(w, "Path variable 'serviceCode' is required", http.StatusBadRequest)
		return
	}

	var req model.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		http.Error(w, "X-Tenant-Id header is required", http.StatusBadRequest)
		return
	}

	if req.Application.TenantId == "" {
		req.Application.TenantId = tenantID
	}
	if req.Application.ServiceCode == "" {
		req.Application.ServiceCode = serviceCode
	}

	ctx := context.Background()
	res, err := c.service.CreateApplication(ctx, req, serviceCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *ApplicationController) SearchApplicationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract path parameter 'serviceCode'
	var criteria model.SearchRequest
	serviceCode := mux.Vars(r)["serviceCode"]

	// Check if serviceCode is missing
	if serviceCode == "" {
		http.Error(w, "'serviceCode' path parameter is required", http.StatusBadRequest)
		return
	}

	// Check if X-Tenant-Id header is present
	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		http.Error(w, "X-Tenant-Id header is required", http.StatusBadRequest)
		return
	}

	// Decode request body

	if criteria.SearchCriteria.TenantId == "" {
		criteria.SearchCriteria.TenantId = tenantID
	}
	if criteria.SearchCriteria.ServiceCode == "" {
		criteria.SearchCriteria.ServiceCode = serviceCode
	}

	module := r.URL.Query().Get("module")
	businessService := r.URL.Query().Get("businessService")
	status := r.URL.Query().Get("status")
	if businessService != "" {
		criteria.SearchCriteria.BusinessService = businessService
	}

	if status != "" {
		criteria.SearchCriteria.Status = status
	}
	if module != "" {
		criteria.SearchCriteria.Module = module
	}
	if idsParam := r.URL.Query().Get("ids"); idsParam != "" {
		criteria.SearchCriteria.Ids = strings.Split(idsParam, ",")
	}
	log.Println("inside search", criteria.SearchCriteria)
	// Prepare context
	ctx := context.Background()

	// Call service to search applications based on criteria
	res, err := c.service.SearchApplication(ctx, criteria.SearchCriteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set response headers and encode the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *ApplicationController) UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract path parameters
	serviceCode := mux.Vars(r)["serviceCode"]
	applicationId := mux.Vars(r)["applicationId"]

	// Check if serviceCode and applicationId are missing
	if serviceCode == "" || applicationId == "" {
		http.Error(w, "Both 'serviceCode' and 'applicationId' path parameters are required", http.StatusBadRequest)
		return
	}

	// Check if X-Tenant-Id header is present
	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		http.Error(w, "X-Tenant-Id header is required", http.StatusBadRequest)
		return
	}

	// Decode request body
	var req model.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Update Service error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println("inside update", req)

	if req.Application.TenantId == "" {
		req.Application.TenantId = tenantID
	}
	if req.Application.ServiceCode == "" {
		req.Application.ServiceCode = serviceCode
	}
	if req.Application.Id == uuid.Nil {
		parsedID, err := uuid.Parse(applicationId)
		if err != nil {
			http.Error(w, "invalid applicationId format:", http.StatusBadRequest)
		}
		req.Application.Id = parsedID
	}
	// Prepare context
	ctx := context.Background()

	// Call service to update the application
	res, err := c.service.UpdateApplication(ctx, req, serviceCode, applicationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
