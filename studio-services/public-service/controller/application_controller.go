package controller

import (
	"context"
	//"crypto/internal/fips140/edwards25519/field"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "os"
	"public-service/model"
	"public-service/service"
	"public-service/utils"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/Priyansuvaish/digit_client/models"
	"github.com/Priyansuvaish/digit_client/services"
	"github.com/xeipuuv/gojsonschema"
)

type ApplicationController struct {
	service            *service.ApplicationService
	workflowIntegrator *service.WorkflowIntegrator
	individualService  *service.IndividualService
	enrichmentService  *service.EnrichmentService
}

func NewApplicationController(service *service.ApplicationService, workflowIntegrator *service.WorkflowIntegrator, individualService *service.IndividualService, enrichmentService *service.EnrichmentService) *ApplicationController {
	return &ApplicationController{service: service, workflowIntegrator: workflowIntegrator, individualService: individualService, enrichmentService: enrichmentService}
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
	filter := map[string]interface{}{
		"module":  req.Application.Module,
		"service": req.Application.BusinessService,
	}
	MdmsCriteria := models.MdmsCriteriaV2Builder().WithTenantID(tenantID).WithFilterMap(filter).WithSchemaCode("Studio.ServiceConfiguration").WithLimit(10).WithOffeset(0).Build()
	mdmsserice := services.NewMdmsV2Service(nil, "egov-mdms-service/v2")
	requestInfo := &models.RequestInfo{
		APIID:     req.RequestInfo.ApiId,
		Ver:       req.RequestInfo.Ver,
		Ts:        int64(req.RequestInfo.Ts),
		Action:    req.RequestInfo.Action,
		Did:       req.RequestInfo.Did,
		Key:       req.RequestInfo.Key,
		MsgID:     req.RequestInfo.MsgId,
		AuthToken: req.RequestInfo.AuthToken,
		UserInfo: map[string]interface{}{
			"uuid":          req.RequestInfo.UserInfo.Uuid,
			"userName":      req.RequestInfo.UserInfo.UserName,
			"name":          req.RequestInfo.UserInfo.Name,
			"emailId":       req.RequestInfo.UserInfo.EmailId,
			"mobileNumber":  req.RequestInfo.UserInfo.MobileNumber,
			"roles":         req.RequestInfo.UserInfo.Roles,
			"tenantId":      req.RequestInfo.UserInfo.TenantId,
			"locale":        req.RequestInfo.UserInfo.Locale,
			"type":          req.RequestInfo.UserInfo.Type,
			"active":        req.RequestInfo.UserInfo.Active,
			"permanentCity": req.RequestInfo.UserInfo.PermanentCity,
		},
	}
	// fmt.Println("requestinfo", requestInfo.ToMap())
	// fmt.Println("mdmscreatia", MdmsCriteria.ToMap())
	search, err := mdmsserice.SearchMDMS(MdmsCriteria, requestInfo)
	//ServiceConfiguration, ok := search.(map[string]interface{})["ServiceConfiguration"]
	if err != nil {
		log.Printf("API call failed: %v", err)
	}

	searchMap, ok := search.(map[string]interface{})
	if !ok {
		log.Println("Failed to assert search as map[string]interface{}")
		return
	}
	mdmsSlice, ok := searchMap["mdms"].([]interface{})
	if !ok || len(mdmsSlice) == 0 {
		log.Println("No mdms data found")
		return
	}

	firstMdmsEntry, ok := mdmsSlice[0].(map[string]interface{})
	if !ok {
		log.Println("Invalid mdms entry format")
		return
	}

	dataMap, ok := firstMdmsEntry["data"].(map[string]interface{})
	if !ok {
		log.Println("No data field found in mdms entry")
		return
	}

	fields, ok := dataMap["fields"]
	if !ok {
		log.Println("No fields found in data")
		return
	}

	// Print only the servicedetail part
	// fmt.Println("servicedetail", req.Application.ServiceDetails)

	// Print only the fields part
	// fmt.Printf("Fields: %+v\n", fields)

	// Validate the service details against the fields schema
	if err := validateServiceDetailsWithSchema(req.Application.ServiceDetails, fields); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Service details validation failed: "+err.Error())
		return
	}

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
	ctx := context.Background()
	log.Println("inside CreateApplicationHandler")
	res, err := c.service.CreateApplication(ctx, req, serviceCode)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Call workflow integrator on success
	err = c.workflowIntegrator.CallWorkflow(&res, req)
	if err != nil {
		log.Printf("Workflow integration failed: %v", err)
		// Optional: return HTTP error or log only
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
	for i := range res.Application {
		err = c.workflowIntegrator.SearchWorkflow(&res.Application[i], criteria.RequestInfo)
		if err != nil {
			log.Printf("Workflow integration failed for application %s: %v", res.Application[i].Id, err)
			// Optional: handle error per item or break early
		}
	}
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
	if req.Application.Workflow.Action == "PAY" {

	}
	res, err := c.service.UpdateApplication(ctx, req, serviceCode, applicationId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Call workflow integrator on success
	err = c.workflowIntegrator.CallWorkflow(&res, req)
	if err != nil {
		log.Printf("Workflow integration failed: %v", err)
		// Optional: return HTTP error or log only
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func ValidateTradeName(tradeName string) error {
	// Check required
	if tradeName == "" {
		return errors.New("tradeName is required")
	}

	// Check length
	length := utf8.RuneCountInString(tradeName)
	if length < 2 {
		return errors.New("tradeName must be at least 2 characters")
	}
	if length > 128 {
		return errors.New("tradeName must be at most 128 characters")
	}

	// Check regex: only letters, digits, spaces
	match, err := regexp.MatchString(`^[A-Za-z0-9 ]+$`, tradeName)
	if err != nil {
		return errors.New("regex match error")
	}
	if !match {
		return errors.New("Only letters and numbers allowed in tradeName")
	}

	return nil
}

// validateServiceDetailsWithSchema validates the service details against the fields schema
func validateServiceDetailsWithSchema(serviceDetails interface{}, fieldsSchema interface{}) error {
	// Convert the fields schema to a proper JSON Schema
	jsonSchema, err := convertFieldsToJSONSchema(fieldsSchema)
	if err != nil {
		return fmt.Errorf("failed to convert fields to JSON schema: %v", err)
	}

	// Convert the schema to JSON string
	schemaBytes, err := json.Marshal(jsonSchema)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON schema: %v", err)
	}

	// Print the generated schema for debugging
	fmt.Printf("Generated JSON Schema: %s\n", string(schemaBytes))

	// Convert the service details to JSON string
	detailsBytes, err := json.Marshal(serviceDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal service details: %v", err)
	}

	// Create schema and document loaders
	schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))
	documentLoader := gojsonschema.NewStringLoader(string(detailsBytes))

	// Validate the document against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Check if the validation was successful
	if !result.Valid() {
		var errorMessages []string
		for _, err := range result.Errors() {
			errorMessages = append(errorMessages, fmt.Sprintf("- %s", err))
		}
		return fmt.Errorf("service details validation failed:\n%s", strings.Join(errorMessages, "\n"))
	}

	return nil
}

// convertFieldsToJSONSchema converts the nested fields structure to a proper JSON Schema
func convertFieldsToJSONSchema(fieldsSchema interface{}) (map[string]interface{}, error) {
	// Convert fields to a slice if it's not already
	fieldsSlice, ok := fieldsSchema.([]interface{})
	if !ok {
		return nil, fmt.Errorf("fields schema is not a slice")
	}

	// Create the root schema object
	schema := map[string]interface{}{
		"$schema":              "http://json-schema.org/draft-07/schema#",
		"type":                 "object",
		"properties":           map[string]interface{}{},
		"additionalProperties": false,
	}

	rootProperties := schema["properties"].(map[string]interface{})

	// Process each top-level field section in the fields schema
	for _, sectionField := range fieldsSlice {
		sectionMap, ok := sectionField.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract section name and type
		sectionName, hasName := sectionMap["name"].(string)
		if !hasName {
			continue
		}

		sectionType, hasType := sectionMap["type"].(string)
		if !hasType {
			continue
		}

		// Handle different section types
		switch sectionType {
		case "object":
			// Create object schema for this section
			sectionSchema := map[string]interface{}{
				"type":                 "object",
				"properties":           map[string]interface{}{},
				"additionalProperties": false,
			}

			// Process properties within this section
			if properties, hasProps := sectionMap["properties"].([]interface{}); hasProps {
				sectionProperties := sectionSchema["properties"].(map[string]interface{})
				requiredProps := []string{}

				for _, prop := range properties {
					propMap, ok := prop.(map[string]interface{})
					if !ok {
						continue
					}

					propName, hasName := propMap["name"].(string)
					if !hasName {
						continue
					}

					// Skip properties with reference
					if _, hasRef := propMap["reference"]; hasRef {
						continue
					}

					// Create property schema
					propSchema := createPropertySchema(propMap)

					// Add to required list if needed
					if isRequired, ok := propMap["required"].(bool); ok && isRequired {
						requiredProps = append(requiredProps, propName)
					}

					// Add property to section schema
					sectionProperties[propName] = propSchema
				}

				// Add required properties if any
				if len(requiredProps) > 0 {
					sectionSchema["required"] = requiredProps
				}
			}

			// Add section to root schema
			rootProperties[sectionName] = sectionSchema

		case "array":
			// Create array schema for this section
			arraySchema := map[string]interface{}{
				"type": "array",
			}

			// Process items schema if available
			if items, hasItems := sectionMap["items"].(map[string]interface{}); hasItems {
				itemsSchema := map[string]interface{}{
					"type":                 "object",
					"properties":           map[string]interface{}{},
					"additionalProperties": false,
				}

				// Process properties within items
				if itemProps, hasItemProps := items["properties"].([]interface{}); hasItemProps {
					itemProperties := itemsSchema["properties"].(map[string]interface{})
					requiredProps := []string{}

					for _, prop := range itemProps {
						propMap, ok := prop.(map[string]interface{})
						if !ok {
							continue
						}

						propName, hasName := propMap["name"].(string)
						if !hasName {
							continue
						}

						// Skip properties with reference
						if _, hasRef := propMap["reference"]; hasRef {
							continue
						}

						// Create property schema
						propSchema := createPropertySchema(propMap)

						// Add to required list if needed
						if isRequired, ok := propMap["required"].(bool); ok && isRequired {
							requiredProps = append(requiredProps, propName)
						}

						// Add property to items schema
						itemProperties[propName] = propSchema
					}

					// Add required properties if any
					if len(requiredProps) > 0 {
						itemsSchema["required"] = requiredProps
					}
				}

				// Add items schema to array schema
				arraySchema["items"] = itemsSchema
			}

			// Add section to root schema
			rootProperties[sectionName] = arraySchema
		}
	}

	return schema, nil
}

// createPropertySchema creates a JSON schema for a single property
func createPropertySchema(propMap map[string]interface{}) map[string]interface{} {
	propSchema := map[string]interface{}{}

	// Set field type
	if propType, ok := propMap["type"].(string); ok {
		switch propType {
		case "string", "number", "integer", "boolean", "array", "object":
			propSchema["type"] = propType
		case "date":
			// For date type, use string with date format
			propSchema["type"] = "string"
			propSchema["format"] = "date"
		default:
			propSchema["type"] = "string" // Default to string
		}
	} else {
		propSchema["type"] = "string" // Default type
	}

	// Set format if available and not already set
	if _, hasFormat := propSchema["format"]; !hasFormat {
		if format, ok := propMap["format"].(string); ok {
			// Handle special formats
			switch format {
			case "radioordropdown":
				// No special format for these UI controls
			case "text":
				// No special format for text
			default:
				propSchema["format"] = format
			}
		}
	}

	// Set min/max length for strings
	if propSchema["type"] == "string" {
		if minLength, ok := propMap["minLength"].(float64); ok {
			propSchema["minLength"] = int(minLength)
		}

		if maxLength, ok := propMap["maxLength"].(float64); ok {
			propSchema["maxLength"] = int(maxLength)
		}
	}

	// Set pattern from validation.regex if available
	if validation, ok := propMap["validation"].(map[string]interface{}); ok {
		if regex, ok := validation["regex"].(string); ok {
			propSchema["pattern"] = regex
		}
	}

	// Set default value if available
	if defaultValue, ok := propMap["defaultValue"]; ok {
		propSchema["default"] = defaultValue
	}

	// Handle dependencies if needed
	if dependencies, ok := propMap["dependencies"].([]interface{}); ok && len(dependencies) > 0 {
		// In JSON Schema, dependencies are more complex
		// For simplicity, we'll just note that this property has dependencies
		propSchema["description"] = "This field has dependencies"
	}

	return propSchema
}
