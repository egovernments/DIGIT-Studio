package service

import (
	"encoding/json"
	"log"
	"public-service/model"
	"public-service/model/individual"
	"strconv"
)

type EnrichmentService struct {
	individualService *IndividualService
}

func NewEnrichmentService(individualService *IndividualService) *EnrichmentService {
	return &EnrichmentService{individualService: individualService}
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
