package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"public-service/config"
	"public-service/model"
)

const (
	TENANT_ID_KEY        = "tenantId"
	BUSINESS_SERVICE_KEY = "businessService"
	ACTION_KEY           = "action"
	COMMENT_KEY          = "comment"
	MODULE_NAME_KEY      = "moduleName"
	BUSINESS_ID_KEY      = "businessId"
	DOCUMENTS_KEY        = "documents"
	ASSIGNEE_KEY         = "assignes"
	UUID_KEY             = "uuid"
	WORKFLOW_REQUEST_KEY = "ProcessInstances"
	REQUEST_INFO_KEY     = "RequestInfo"
)

// WorkflowIntegrator handles interaction with the workflow service.
type WorkflowIntegrator struct {
	HttpClient *http.Client
}

// NewWorkflowIntegrator returns a new instance of WorkflowIntegrator.
func NewWorkflowIntegrator() *WorkflowIntegrator {
	return &WorkflowIntegrator{
		HttpClient: &http.Client{},
	}
}

// CallWorkflow integrates with the workflow and updates the application with the workflow response.
func (wi *WorkflowIntegrator) CallWorkflow(applicationResponse *model.ApplicationResponse, req model.ApplicationRequest) error {
	app := applicationResponse.Application
	log.Println("inside CallWorkflow")
	log.Println("🔥🔥🔥 Inside CallWorkflow - LOG TRIGGERED 🔥🔥🔥")
	requestPayload := make(map[string]interface{})
	processInstance := make(map[string]interface{})

	processInstance[BUSINESS_ID_KEY] = app.ApplicationNumber
	processInstance[TENANT_ID_KEY] = app.TenantId
	processInstance[BUSINESS_SERVICE_KEY] = app.BusinessService
	processInstance[MODULE_NAME_KEY] = app.Module
	processInstance[ACTION_KEY] = app.Workflow.Action
	log.Println("workflow")
	processInstance[COMMENT_KEY] = app.Workflow.Comment

	if len(app.Workflow.Assignees) > 0 {
		var uuidMaps []map[string]string
		for _, assignee := range app.Workflow.Assignees {
			uuidMaps = append(uuidMaps, map[string]string{UUID_KEY: assignee.Uuid.String()})
		}
		processInstance[ASSIGNEE_KEY] = uuidMaps
	}

	processInstance[DOCUMENTS_KEY] = app.Workflow.Documents

	requestPayload[REQUEST_INFO_KEY] = req.RequestInfo
	requestPayload[WORKFLOW_REQUEST_KEY] = []map[string]interface{}{processInstance}
	log.Println("inside CallWorkflow", requestPayload)
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return fmt.Errorf("error marshaling workflow request: %w", err)
	}

	wfHost := os.Getenv("WORKFLOW_HOST")
	wfPath := os.Getenv("WORKFLOW_TRANSITION_PATH")
	if wfHost == "" || wfPath == "" {
		log.Println("wfHost", wfHost)
		log.Println("wfPath", wfPath)
		return errors.New("workflow host or path is not set in environment variables")
	}

	url := wfHost + wfPath
	resp, err := wi.HttpClient.Post(url, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to call workflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("workflow service returned status %d: %v", resp.StatusCode, errResp)
	}

	var wfResponse model.ProcessInstanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&wfResponse); err != nil {
		return fmt.Errorf("error decoding workflow response: %w", err)
	}

	if len(wfResponse.ProcessInstances) == 0 {
		return errors.New("no process instance returned from workflow")
	}
	app.ProcessInstance = &wfResponse.ProcessInstances
	applicationResponse.Application.ProcessInstance = &wfResponse.ProcessInstances
	return nil
}

func (wi *WorkflowIntegrator) SearchWorkflow(applicationResponse *model.Application, req model.RequestInfo) error {
	app := applicationResponse
	log.Println("Search CallWorkflow")
	log.Println("🔥🔥🔥 Inside SearchWorkflow - LOG TRIGGERED 🔥🔥🔥")
	requestPayload := make(map[string]interface{})
	requestPayload[REQUEST_INFO_KEY] = req

	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return fmt.Errorf("error marshaling workflow request: %w", err)
	}

	config.LoadEnv()
	wfHost := config.GetEnv("WORKFLOW_HOST")
	wfPath := config.GetEnv("WORKFLOW_SEARCH_PATH")
	if wfHost == "" || wfPath == "" {
		log.Println("wfHost", wfHost)
		log.Println("wfPath", wfPath)
		return errors.New("workflow host or search path is not set in environment variables")
	}

	url := fmt.Sprintf("%s%s?tenantId=%s&businessIds=%s", wfHost, wfPath, app.TenantId, app.ApplicationNumber)
	log.Println("URL:", url)

	resp, err := wi.HttpClient.Post(url, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to call workflow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("workflow service returned status %d: %v", resp.StatusCode, errResp)
	}

	var wfResponse model.ProcessInstanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&wfResponse); err != nil {
		return fmt.Errorf("error decoding workflow response: %w", err)
	}

	if len(wfResponse.ProcessInstances) == 0 {
		return errors.New("no process instance returned from workflow")
	}
	applicationResponse.ProcessInstance = &wfResponse.ProcessInstances
	return nil
}
