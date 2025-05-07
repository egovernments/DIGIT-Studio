package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
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
}

func NewEnrichmentService(individualService *IndividualService, demandService *DemandService, mdmsService *MDMSService) *EnrichmentService {
	return &EnrichmentService{individualService: individualService, DemandService: demandService, MDMSService: mdmsService}
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

func (s *EnrichmentService) EnrichApplicationsWithDemand(apps model.ApplicationRequest) model.ApplicationRequest {
	if apps.Application.Workflow.Action == "APPLY" {

		// Create MasterDetail
		masterDetail := model.MasterDetail{
			Name:   os.Getenv("SERVICE_MASTER_NAME"),
			Filter: "", // Add filter if needed
		}

		// Create ModuleDetail
		moduleDetail := model.ModuleDetail{
			ModuleName:    os.Getenv("SERVICE_MODULE_NAME"),
			MasterDetails: []model.MasterDetail{masterDetail},
		}

		// Create MdmsCriteria
		criteria := model.MdmsCriteria{
			TenantID:      apps.Application.TenantId,
			ModuleDetails: []model.ModuleDetail{moduleDetail},
		}

		// Call MDMS search
		mdmsResponse, err := s.MDMSService.MDMSSearch(criteria, apps.RequestInfo)
		if err != nil {
			log.Printf("Failed to fetch MDMS data: %v", err)
			return apps
		}
		log.Println("MDMS Response:", mdmsResponse)
		_ = demand.Demand{
			ID:              uuid.NewString(),
			TenantID:        apps.Application.TenantId,
			ConsumerCode:    apps.Application.ApplicationNumber,
			ConsumerType:    apps.Application.Module,
			BusinessService: apps.Application.BusinessService,
			Payer:           apps.RequestInfo.UserInfo,
		}

		return apps
	}
	return apps
}
