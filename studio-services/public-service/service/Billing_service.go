package service

import (
	"fmt"
	"log"
	"public-service/config"
	"public-service/model"
	"public-service/model/demand"
	"public-service/repository"
)

type DemandService struct {
	restCallRepo repository.RestCallRepository
}

func NewDemandService(repo repository.RestCallRepository) *DemandService {
	return &DemandService{
		restCallRepo: repo,
	}
}

// SaveDemand creates demand records via external billing service
func (r *DemandService) SaveDemand(requestInfo model.RequestInfo, demands []demand.Demand) ([]demand.Demand, error) {
	//url := fmt.Sprintf("%s%s",os.Getenv(), Config.DemandCreateEndPoint)
	url := config.GetEnv("BILLING_SERVICE_HOST") + config.GetEnv("DEMAND_CREATE_ENDPOINT")
	demandRequest := demand.DemandRequest{
		RequestInfo: requestInfo,
		Demands:     demands,
	}

	log.Printf("Creating demand for consumer code: %s", demandRequest.Demands[0].ConsumerCode)

	//responseBytes, err := r.ServiceRequestRepo.FetchResult(url, demandRequest)
	resp := demand.DemandResponse{}
	err := r.restCallRepo.Post(url, demandRequest, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call billing service: %w", err)
	}

	return resp.Demands, nil
}

// UpdateDemand updates demand records via external billing service
func (r *DemandService) UpdateDemand(requestInfo model.RequestInfo, demands []demand.Demand) ([]demand.Demand, error) {

	url := config.GetEnv("BILLING_SERVICE_HOST") + config.GetEnv("DEMAND_UPDATE_ENDPOINT")

	demandRequest := demand.DemandRequest{
		RequestInfo: requestInfo,
		Demands:     demands,
	}

	log.Printf("Updating demand for consumer code: %s", demandRequest.Demands[0].ConsumerCode)

	resp := demand.DemandResponse{}
	err := r.restCallRepo.Post(url, demandRequest, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call billing service: %w", err)
	}

	return resp.Demands, nil

}

// GetDemandsToAddPenalty returns a list of consumer codes eligible for penalty
/*func (r *DemandRepository) GetDemandsToAddPenalty(tenantId string, penaltyThresholdTime *big.Int, penaltyApplicableAfterDays int) ([]string, error) {
	query, args := r.QueryBuilder.GetPenaltyQuery(tenantId, penaltyThresholdTime, penaltyApplicableAfterDays)
	log.Printf("query: %s", query)

	var consumerCodes []string
	if err := r.DB.Select(&consumerCodes, query, args...); err != nil {
		return nil, fmt.Errorf("error querying demands for penalty: %w", err)
	}

	return consumerCodes, nil
}
*/
