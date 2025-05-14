package service

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"public-service/model"
	"public-service/model/demand"
	"public-service/model/individual"
	"strconv"

	"github.com/google/uuid"
)

type EnrichmentService struct {
	individualService *IndividualService
	DemandService     *DemandService
	MDMSService       *MDMSService
	MDMSV2Service     *MDMSV2Service
	IdGenService      *IdGenService
	SMSService        *SMSService
}

func NewEnrichmentService(individualService *IndividualService, demandService *DemandService, mdmsService *MDMSService, mdmsServiceV2 *MDMSV2Service, idGenService *IdGenService, smsService *SMSService) *EnrichmentService {
	return &EnrichmentService{individualService: individualService, DemandService: demandService, MDMSService: mdmsService, MDMSV2Service: mdmsServiceV2, IdGenService: idGenService}
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
	if apps.Application.Workflow.Action == "VERIFY_AND_FORWARD" {
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
		var payerUser model.User // Replace with your actual User struct type
		if len(apps.Application.Applicants) > 0 {
			individualId := apps.Application.Applicants[0].UserId // Assuming this exists

			// Fetch individual/user details using the ID
			criteria := map[string]interface{}{
				"uuid":     individualId,
				"tenantId": apps.Application.TenantId,
			}
			indResp := s.individualService.GetIndividual(model.RequestInfo{}, criteria)
			if jsonBytes, err := json.MarshalIndent(indResp, "", "  "); err == nil {
				log.Printf("Indiviual response:\n%s\n", string(jsonBytes))
			} else {
				log.Printf("Indiviual response (raw): %+v\n", indResp)
			}

			if len(indResp.Individual) > 0 {
				// Map individual data to User (Payer)
				individual := indResp.Individual[0]

				parsedUUID, err := uuid.Parse(individual.UserUuid)
				if err != nil {
					log.Printf("Invalid UUID format for UserUuid: %v", err)
					return apps
				}
				payerUser = model.User{
					Uuid:         parsedUUID,
					UserName:     individual.UserDetails.UserName,
					Name:         individual.Name.GivenName,
					MobileNumber: individual.MobileNumber,
					EmailId:      individual.Email,
					TenantId:     individual.TenantId,
					Type:         individual.UserDetails.Type,
				}
			}
		} else {
			log.Println("No applicants found to assign as payer")
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
			ConsumerCode:    apps.Application.ApplicationNumber,
			ConsumerType:    apps.Application.Module,
			BusinessService: businessService,
			Payer:           &payerUser,
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
			_, err2 := s.SMSService.SendSMS(apps, apps.Application.TenantId, "DIGIT_STUDIO_DEMAND_CREATED", apps.Application.Applicants)
			if err2 != nil {
				log.Printf("Failed to send SMS: %v", err2)
			}

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

func (s *EnrichmentService) EnrichApplicationsWithIdGen(apps model.ApplicationRequest) model.ApplicationRequest {
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

	var format string
	var name string
	firstEntry, ok := mdmsList[0].(map[string]interface{})
	if !ok {
		log.Println("Invalid MDMS format: first entry is not a map")
		return apps
	}

	data, ok := firstEntry["data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid MDMS format: missing or invalid 'data'")
		return apps
	}

	idGens, ok := data["idgen"].([]interface{})
	if !ok || len(idGens) == 0 {
		log.Println("No 'idgen' section in MDMS data")
	}

	for _, item := range idGens {
		idGen, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if idGenType, _ := idGen["Type"].(string); idGenType == "application" {
			format, _ = idGen["format"].(string)
			name, _ = idGen["name"].(string)
			break
		}
	}

	// Validate if name and format were found
	if name == "" || format == "" {
		log.Println("IDGen config not found for application type")
		name = "public-service.application.id"
		format = "APL-[cy:yyyy-MM-dd]-[SEQ_PUBLIC_APPLICATION]"

	}
	name = "public-service.application.id"
	format = "APL-[cy:yyyy-MM-dd]-[SEQ_PUBLIC_APPLICATION]"

	// Count should be at least 1
	ids, err := s.IdGenService.GetId(apps.RequestInfo, apps.Application.TenantId, name, format, 1)
	if err != nil {
		log.Printf("Error getting ID from IDGenService: %v", err)
		return apps
	}
	if len(ids) > 0 {
		apps.Application.ApplicationNumber = ids[0]
	}

	return apps
}
