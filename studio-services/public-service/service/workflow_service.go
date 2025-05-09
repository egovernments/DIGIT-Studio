package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"public-service/model"
)

type WorkflowService struct {
	httpClient *http.Client
}

func NewWorkflowService() *WorkflowService {
	return &WorkflowService{
		httpClient: &http.Client{},
	}
}

func (ws *WorkflowService) GetBusinessService(serviceRequest model.ServiceRequest, requestInfo model.RequestInfo, applicationNumber string) (*model.BusinessService, error) {
	url := ws.buildSearchURL(serviceRequest, true, applicationNumber)

	workflowReq := model.RequestInfoWrapper{
		RequestInfo: requestInfo,
	}

	payload, err := json.Marshal(workflowReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request info: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch business service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from business service fetch: %d", resp.StatusCode)
	}

	var response model.BusinessServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.BusinessServices) == 0 {
		return nil, errors.New("no business services found")
	}

	return &response.BusinessServices[0], nil
}

func (ws *WorkflowService) IsStateUpdatable(status string, businessService *model.BusinessService) bool {
	for _, state := range businessService.States {
		if state.ApplicationStatus != nil && *state.ApplicationStatus == status {
			return state.IsStateUpdatable
		}
	}
	return false
}

func (ws *WorkflowService) GetCurrentState(status string, businessService *model.BusinessService) *string {
	for _, state := range businessService.States {
		if state.ApplicationStatus != nil && *state.ApplicationStatus == status {
			return &state.State
		}
	}
	return nil
}

func (ws *WorkflowService) GetCurrentStateObj(status string, businessService *model.BusinessService) *model.State {
	for _, state := range businessService.States {
		if state.ApplicationStatus != nil && *state.ApplicationStatus == status {
			return &state
		}
	}
	return nil
}

func (ws *WorkflowService) buildSearchURL(serviceRequest model.ServiceRequest, isBusinessService bool, applicationNumber string) string {
	host := os.Getenv("WF_HOST")
	businessPath := os.Getenv("WF_BUSINESS_SERVICE_SEARCH_PATH")
	processPath := os.Getenv("WF_PROCESS_PATH")

	var url string
	if isBusinessService {
		url = fmt.Sprintf("%s%s?tenantId=%s&businessServices=%s", host, businessPath, serviceRequest.Service.TenantId, serviceRequest.Service.BusinessService)
	} else {
		url = fmt.Sprintf("%s%s?tenantId=%s&businessIds=%s", host, processPath, serviceRequest.Service.TenantId, applicationNumber)
	}
	return url
}

func (ws *WorkflowService) CreateBusinessService(businessServiceRequest model.BusinessServiceRequest) error {
	payload, err := json.Marshal(businessServiceRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	url := os.Getenv("WF_BUSINESS_SERVICE_CREATE_URL")
	if url == "" {
		return errors.New("workflow business service create URL not set in environment variables")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call business service create API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("business service create returned unexpected status: %d", resp.StatusCode)
	}

	return nil
}
