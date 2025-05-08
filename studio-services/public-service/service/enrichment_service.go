package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"math/big"
	"os"
	"public-service/model"
	"public-service/model/demand"
	"public-service/model/individual"
	"strconv"
)

type EnrichmentService struct {
	individualService *IndividualService
	DemandService     *DemandService
	MDMSService       *MDMSService
	MDMSV2Service     *MDMSV2Service
}

func NewEnrichmentService(individualService *IndividualService, demandService *DemandService, mdmsService *MDMSService, mdmsServiceV2 *MDMSV2Service) *EnrichmentService {
	return &EnrichmentService{individualService: individualService, DemandService: demandService, MDMSService: mdmsService, MDMSV2Service: mdmsServiceV2}
}

func (s *EnrichmentService) EnrichApplicationsWithIndividuals(apps []model.Application, criteria model.SearchCriteria) []model.Application {
	userCache := make(map[string]individual.Individual)

	for aIndex, app := range apps {
		for i, applicant := range app.Applicants {
			if applicant.UserId == "" {
				continue
			}
			if cached, ok := userCache[applicant.UserId]; ok {
				app.Applicants[i].Name = cached.Name.GivenName
				app.Applicants[i].MobileNumber, _ = strconv.ParseInt(cached.MobileNumber, 10, 64)
				app.Applicants[i].EmailId = cached.Email
				continue
			}
			criteria := map[string]interface{}{
				"uuid":     applicant.UserId,
				"tenantId": criteria.TenantId,
			}
			indResp := s.individualService.GetIndividual(model.RequestInfo{}, criteria)
			if jsonBytes, err := json.MarshalIndent(indResp, "", "  "); err == nil {
				log.Printf("Indiviual response:\n%s\n", string(jsonBytes))
			} else {
				log.Printf("Indiviual response (raw): %+v\n", indResp)
			}
			if len(indResp.Individual) > 0 {
				ind := indResp.Individual[0]
				userCache[applicant.UserId] = ind
				app.Applicants[i].Name = ind.Name.GivenName
				app.Applicants[i].MobileNumber, _ = strconv.ParseInt(ind.MobileNumber, 10, 64)
				app.Applicants[i].EmailId = ind.Email
			}
		}
		apps[aIndex] = app
	}

	return apps
}

/*func (s *EnrichmentService) EnrichApplicationsWithDemand(apps model.ApplicationRequest) model.ApplicationRequest {
if apps.Application.Workflow.Action == "PAY" {

	schemaCode := os.Getenv("SERVICE_MODULE_NAME") + "." + os.Getenv("SERVICE_MASTER_NAME")
	mdmsData, _ := s.MDMSV2Service.SearchMDMS(apps.Application.TenantId, schemaCode, apps.Application.BusinessService, apps.Application.Module, apps.RequestInfo)
	mdmsList, ok := mdmsData["mdms"].([]interface{})
	if !ok || len(mdmsList) == 0 {
		log.Println("MDMS data missing or invalid")
		return apps
	}

	firstEntry, ok := mdmsList[0].(map[string]interface{})
	if !ok {
		log.Println("Invalid MDMS entry format")
		return apps
	}

	data, ok := firstEntry["data"].(map[string]interface{})
	if !ok {
		log.Println("MDMS 'data' field missing or invalid")
		return apps
	}

	billData, ok := data["bill"].(map[string]interface{})
	if !ok {
		log.Println("No 'bill' section in MDMS data")
		return apps
	}

	serviceName, ok := billData["service"].(string)
	if !ok {
		log.Println("No 'service' field inside bill section")
		return apps
	}

	log.Println("Extracted Bill Service:", serviceName)
	// TODO: Replace with actual logic to extract the TaxHeadMasterCode from mdmsResponse
	taxHeadCode := "TL_APPLICATION_FEE"*/
/*	schemaCode := os.Getenv("SERVICE_MODULE_NAME") + "." + os.Getenv("SERVICE_MASTER_NAME")
		mdmsData, _ := s.MDMSV2Service.SearchMDMS(
			apps.Application.TenantId,
			schemaCode,
			apps.Application.BusinessService,
			apps.Application.Module,
			apps.RequestInfo,
		)

		mdmsList, ok := mdmsData["mdms"].([]interface{})
		if !ok || len(mdmsList) == 0 {
			log.Println("MDMS data missing or invalid")
			return apps
		}

		firstEntry, _ := mdmsList[0].(map[string]interface{})
		data, _ := firstEntry["data"].(map[string]interface{})
		billData, ok := data["bill"].(map[string]interface{})
		if !ok {
			log.Println("No 'bill' section in MDMS data")
			return apps
		}

		// Step 2: Extract taxHeadCode from bill.taxHead
		var taxHeadCode string
		if taxHeads, ok := billData["taxHead"].([]interface{}); ok {
			for _, item := range taxHeads {
				taxHead, _ := item.(map[string]interface{})
				if taxHead["service"] == apps.Application.BusinessService && taxHead["code"] == "TL_TAX" {
					taxHeadCode = taxHead["code"].(string)
					break
				}
			}
		}
		if taxHeadCode == "" {
			log.Println("No matching taxHeadCode found, defaulting")
			taxHeadCode = "TL_TAX"
		}

		// Step 3: Extract taxPeriodFrom and taxPeriodTo from bill.taxPeriod
		var taxPeriodFrom, taxPeriodTo *big.Float
		if taxPeriods, ok := billData["taxPeriod"].([]interface{}); ok {
			for _, item := range taxPeriods {
				tp, _ := item.(map[string]interface{})
				if tp["service"] == apps.Application.BusinessService {
					if from, ok := tp["fromDate"].(float64); ok {
						taxPeriodFrom = big.NewFloat(from)
					}
					if to, ok := tp["toDate"].(float64); ok {
						taxPeriodTo = big.NewFloat(to)
					}
					break
				}
			}
		}

		// Step 4: Extract businessService from bill.BusinessService
		var businessService string
		if bsMap, ok := billData["BusinessService"].(map[string]interface{}); ok {
			if bsMap["businessService"] == apps.Application.BusinessService {
				if code, ok := bsMap["code"].(string); ok {
					businessService = code
				}
			}
		}
		if businessService == "" {
			businessService = apps.Application.BusinessService // fallback
		}

		demandDetail := demand.DemandDetail{
			ID:                uuid.NewString(),
			TaxHeadMasterCode: taxHeadCode,
			TaxAmount:         big.NewFloat(50.0),
			CollectionAmount:  big.NewFloat(0.0),
			TenantID:          apps.Application.TenantId,
			AuditDetails:      nil, // Populate if needed
		}

		d := demand.Demand{
			ID:              uuid.NewString(),
			TenantID:        apps.Application.TenantId,
			ConsumerCode:    apps.Application.ApplicationNumber,
			ConsumerType:    apps.Application.Module,
			BusinessService: apps.Application.BusinessService,
			Payer:           apps.RequestInfo.UserInfo,
			TaxPeriodFrom:   taxPeriodFrom,
			TaxPeriodTo:     taxPeriodTo,
			DemandDetails:   []demand.DemandDetail{demandDetail},
			AuditDetails:    nil,
		}

		logJSON("Created Demand", d)

		// Optional: attach the demand to the application if needed
	}

	return apps
}*/

func (s *EnrichmentService) EnrichApplicationsWithDemand(apps model.ApplicationRequest) model.ApplicationRequest {
	if apps.Application.Workflow.Action == "PAY" {
		schemaCode := os.Getenv("SERVICE_MODULE_NAME") + "." + os.Getenv("SERVICE_MASTER_NAME")
		mdmsData, _ := s.MDMSV2Service.SearchMDMS(
			apps.Application.TenantId,
			schemaCode,
			apps.Application.BusinessService,
			apps.Application.Module,
			apps.RequestInfo,
		)

		mdmsList, ok := mdmsData["mdms"].([]interface{})
		if !ok || len(mdmsList) == 0 {
			log.Println("MDMS data missing or invalid")
			return apps
		}

		firstEntry, _ := mdmsList[0].(map[string]interface{})
		data, _ := firstEntry["data"].(map[string]interface{})
		billData, ok := data["bill"].(map[string]interface{})
		if !ok {
			log.Println("No 'bill' section in MDMS data")
			return apps
		}

		// Step 2: Extract taxHeadCode from bill.taxHead
		var taxHeadCode string
		if taxHeads, ok := billData["taxHead"].([]interface{}); ok {
			for _, item := range taxHeads {
				taxHead, _ := item.(map[string]interface{})
				if taxHead["service"] == apps.Application.BusinessService && taxHead["code"] == "TL_TAX" {
					taxHeadCode = taxHead["code"].(string)
					break
				}
			}
		}
		if taxHeadCode == "" {
			log.Println("No matching taxHeadCode found, defaulting")
			taxHeadCode = "TL_TAX"
		}

		// Step 3: Extract taxPeriodFrom and taxPeriodTo from bill.taxPeriod
		var taxPeriodFrom, taxPeriodTo *int64
		if taxPeriods, ok := billData["taxPeriod"].([]interface{}); ok {
			for _, item := range taxPeriods {
				tp, _ := item.(map[string]interface{})
				if tp["service"] == apps.Application.BusinessService {
					if from, ok := tp["fromDate"].(float64); ok {
						fromInt := int64(from)
						taxPeriodFrom = &fromInt
					}
					if to, ok := tp["toDate"].(float64); ok {
						toInt := int64(to)
						taxPeriodTo = &toInt
					}
					break
				}
			}
		}

		// Step 4: Extract businessService from bill.BusinessService
		var businessService string
		if bsMap, ok := billData["BusinessService"].(map[string]interface{}); ok {
			if bsMap["businessService"] == apps.Application.BusinessService {
				if code, ok := bsMap["code"].(string); ok {
					businessService = code
				}
			}
		}
		if businessService == "" {
			businessService = apps.Application.BusinessService // fallback
		}

		demandDetail := demand.DemandDetail{
			ID:                uuid.NewString(),
			TaxHeadMasterCode: taxHeadCode,
			TaxAmount:         big.NewFloat(50.0),
			CollectionAmount:  big.NewFloat(0.0),
			TenantID:          apps.Application.TenantId,
			AuditDetails:      nil,
		}

		d := demand.Demand{
			ID:              uuid.NewString(),
			TenantID:        apps.Application.TenantId,
			ConsumerCode:    "APL-03993",
			ConsumerType:    apps.Application.Module,
			BusinessService: apps.Application.BusinessService,
			Payer:           apps.RequestInfo.UserInfo,
			TaxPeriodFrom:   taxPeriodFrom,
			TaxPeriodTo:     taxPeriodTo,
			DemandDetails:   []demand.DemandDetail{demandDetail},
			AuditDetails:    nil,
		}
		var demands []demand.Demand
		demands = append(demands, d)

		createdDemands, err := s.DemandService.SaveDemand(apps.RequestInfo, demands)
		if err != nil {
			log.Printf("Failed to save demand: %v", err)
		} else {
			logJSON("Saved Demands Response", createdDemands)
		}

		// Optional: attach the demand to the application if needed
	}

	return apps
}

func logJSON(message string, data interface{}) {
	if jsonData, err := json.Marshal(data); err == nil {
		log.Printf(`{"message": "%s", "data": %s}`, message, jsonData)
	} else {
		log.Printf(`{"message": "%s", "error": "%v"}`, message, err)
	}
}

func logError(message string, err error, context interface{}) {
	ctxJSON, _ := json.Marshal(context)
	log.Printf(`{"error": "%s", "details": "%v", "context": %s}`, message, err, ctxJSON)
}
