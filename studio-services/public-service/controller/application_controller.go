package controller

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "os"
	"public-service/model"
	"public-service/service"
	"public-service/utils"
	"strconv"
	"strings"
)

type ApplicationController struct {
	service            *service.ApplicationService
	workflowIntegrator *service.WorkflowIntegrator
	individualService  *service.IndividualService
	enrichmentService  *service.EnrichmentService
	smsService         *service.SMSService
}

func NewApplicationController(service *service.ApplicationService, workflowIntegrator *service.WorkflowIntegrator, individualService *service.IndividualService, enrichmentService *service.EnrichmentService, smsService *service.SMSService) *ApplicationController {
	return &ApplicationController{service: service, workflowIntegrator: workflowIntegrator, individualService: individualService, enrichmentService: enrichmentService, smsService: smsService}
}

func (c *ApplicationController) CreateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	serviceCode := mux.Vars(r)["serviceCode"]

	if serviceCode == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Path variable 'serviceCode' is required")
		return
	}

	var req model.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Missing header 'X-Tenant-Id'")
		return
	}

	if req.Application.TenantId == "" {
		req.Application.TenantId = tenantID
	}
	if req.Application.ServiceCode == "" {
		req.Application.ServiceCode = serviceCode
	}
	req = c.enrichmentService.EnrichApplicationsWithIdGen(req)
	log.Println(req)
	for i, applicant := range req.Application.Applicants {
		criteria := map[string]interface{}{
			"mobileNumber": strconv.FormatInt(applicant.MobileNumber, 10),
			"tenantId":     req.Application.TenantId,
		}

		// Check if individual exists
		resp := c.individualService.GetIndividual(req.RequestInfo, criteria)

		if len(resp.Individual) == 0 {
			// If not found, create individual
			createdResp := c.individualService.CreateUser(applicant, req.RequestInfo)
			if createdResp.Individual.IndividualId != "" {
				req.Application.Applicants[i].UserId = createdResp.Individual.IndividualId
			} else {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create individual")
				return
			}
		} else {
			// Individual exists, update applicant UserId
			req.Application.Applicants[i].UserId = resp.Individual[i].IndividualId
		}
	}
	c.enrichmentService.EnrichApplicationsWithDemand(req)
	// Call workflow integrator on success
	err = c.workflowIntegrator.CallWorkflow(&req)
	if err != nil {
		log.Printf("Workflow integration failed: %v", err)
		// Optional: return HTTP error or log only
	}
	ctx := context.Background()
	log.Println("inside CreateApplicationHandler")
	res, err := c.service.CreateApplication(ctx, req, serviceCode)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err2 := c.smsService.SendSMS(req, req.Application.TenantId, "DIGIT_STUDIO_APPLY_NEW_CONNECTION", req.Application.Applicants)
	if err2 != nil {
		log.Printf("error sending sms ")
	}
	log.Printf("ProcessInstance enriched: %+v", res.Application.ProcessInstance)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *ApplicationController) SearchApplicationHandler(w http.ResponseWriter, r *http.Request) {
	var criteria model.SearchRequest
	serviceCode := mux.Vars(r)["serviceCode"]

	if serviceCode == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Path variable 'serviceCode' is required")
		return
	}

	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Missing header 'X-Tenant-Id'")
		return
	}

	if criteria.SearchCriteria.TenantId == "" {
		criteria.SearchCriteria.TenantId = tenantID
	}
	if criteria.SearchCriteria.ServiceCode == "" {
		criteria.SearchCriteria.ServiceCode = serviceCode
	}

	module := r.URL.Query().Get("module")
	businessService := r.URL.Query().Get("businessService")
	status := r.URL.Query().Get("status")
	applicationNumber := r.URL.Query().Get("applicationNumber")
	if businessService != "" {
		criteria.SearchCriteria.BusinessService = businessService
	}
	if status != "" {
		criteria.SearchCriteria.Status = status
	}
	if module != "" {
		criteria.SearchCriteria.Module = module
	}
	if applicationNumber != "" {
		criteria.SearchCriteria.ApplicationNumber = applicationNumber
	}
	if idsParam := r.URL.Query().Get("ids"); idsParam != "" {
		criteria.SearchCriteria.Ids = strings.Split(idsParam, ",")
	}
	log.Println("inside search", criteria.SearchCriteria)
	ctx := context.Background()
	res, err := c.service.SearchApplication(ctx, criteria.SearchCriteria)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	//TODO: enrich ProcessInstance as request nfo not there its throwing error
	/*for i := range res.Application {
		err = c.workflowIntegrator.SearchWorkflow(&res.Application[i], criteria.RequestInfo)
		if err != nil {
			log.Printf("Workflow integration failed for application %s: %v", res.Application[i].Id, err)
			// Optional: handle error per item or break early
		}
	}*/
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *ApplicationController) UpdateApplicationHandler(w http.ResponseWriter, r *http.Request) {
	serviceCode := mux.Vars(r)["serviceCode"]
	applicationId := mux.Vars(r)["applicationId"]

	if serviceCode == "" || applicationId == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Path variable 'serviceCode' is required")
		return
	}

	tenantID := r.Header.Get("X-Tenant-Id")
	if tenantID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Missing header 'X-Tenant-Id'")
		return
	}

	var req model.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Update Service error: %v", err)
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
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
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid application id")
		}
		req.Application.Id = parsedID
	}
	ctx := context.Background()
	c.enrichmentService.EnrichApplicationsWithDemand(req)
	// Call workflow integrator on success
	err = c.workflowIntegrator.CallWorkflow(&req)
	if err != nil {
		log.Printf("Workflow integration failed: %v", err)
		// Optional: return HTTP error or log only
	}
	res, err := c.service.UpdateApplication(ctx, req, serviceCode, applicationId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("ProcessInstance enriched: %+v", res.Application.ProcessInstance)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
