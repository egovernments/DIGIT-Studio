package controller

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"public-service/model"
	"public-service/service"
)

type PublicController struct {
	service *service.PublicService
}

func NewServiceController(service *service.PublicService) *PublicController {
	return &PublicController{service: service}
}

func (c *PublicController) CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	var req model.ServiceRequest
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

	if req.Service.TenantId == "" {
		req.Service.TenantId = tenantID
	}
	ctx := context.Background()
	res, err := c.service.CreateService(ctx, req, tenantID)
	if err != nil {
		log.Printf("CreateService error: %v", err)
		http.Error(w, "Failed to create service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *PublicController) UpdateServiceHandler(w http.ResponseWriter, r *http.Request) {
	serviceCode := mux.Vars(r)["serviceCode"]
	var serviceRequest model.ServiceRequest

	if err := json.NewDecoder(r.Body).Decode(&serviceRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		http.Error(w, "X-Tenant-Id header is required", http.StatusBadRequest)
		return
	}

	// Ensure the tenantId and serviceCode are set in the request body
	if serviceRequest.Service.TenantId == "" {
		serviceRequest.Service.TenantId = tenantID
	}
	if serviceRequest.Service.ServiceCode == "" {
		serviceRequest.Service.ServiceCode = serviceCode
	}

	ctx := context.Background()
	res, err := c.service.UpdateService(ctx, serviceRequest, serviceCode)
	if err != nil {
		log.Printf("Update Service error: %v", err)
		http.Error(w, "Failed to update service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *PublicController) SearchServiceHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		http.Error(w, "X-Tenant-Id header is required", http.StatusBadRequest)
		return
	}

	module := r.URL.Query().Get("module")
	businessService := r.URL.Query().Get("businessService")
	serviceCode := r.URL.Query().Get("serviceCode")

	searchServiceCriteria := model.SearchCriteria{}
	searchServiceCriteria.TenantId = tenantID

	if businessService != "" {
		searchServiceCriteria.BusinessService = businessService
	}
	if module != "" {
		searchServiceCriteria.Module = module
	}
	if serviceCode != "" {
		searchServiceCriteria.ServiceCode = serviceCode
	}

	ctx := context.Background()
	res, err := c.service.SearchService(ctx, searchServiceCriteria)
	if err != nil {
		log.Printf("SearchService error: %v", err)
		http.Error(w, "Failed to search service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
