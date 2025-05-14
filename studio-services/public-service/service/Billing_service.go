package service

import (
	"encoding/json"
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
	// Pretty-print the request as JSON
	if jsonBytes, err := json.MarshalIndent(demandRequest, "", "  "); err == nil {
		log.Printf("Demand Request JSON:\n%s", string(jsonBytes))
	} else {
		log.Printf("Failed to marshal demandRequest: %v", err)
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

/*
	func (r *DemandRepository) GetDemandsToAddPenalty(tenantId string, penaltyThresholdTime *big.Int, penaltyApplicableAfterDays int) ([]string, error) {
		query, args := r.QueryBuilder.GetPenaltyQuery(tenantId, penaltyThresholdTime, penaltyApplicableAfterDays)
		log.Printf("query: %s", query)

		var consumerCodes []string
		if err := r.DB.Select(&consumerCodes, query, args...); err != nil {
			return nil, fmt.Errorf("error querying demands for penalty: %w", err)
		}

		return consumerCodes, nil
	}
*/
func (r *DemandService) fetchBill(request model.ApplicationRequest) ([]demand.Bill, error) {
	baseURL := config.GetEnv("BILLING_SERVICE_HOST")
	endpoint := config.GetEnv("BILL_FETCH_ENDPOINT") // e.g., /billing-service/bill/v2/_fetchbill
	queryParams := fmt.Sprintf("?tenantId=%s&consumerCode=%s&businessService=%s",
		request.Application.TenantId,
		request.Application.ApplicationNumber,
		request.Application.BusinessService,
	)

	url := baseURL + endpoint + queryParams

	// Prepare request body
	// Have to make IT citizen as type employee Cannot search state level tenantid Bills
	requestInfo := request.RequestInfo
	requestInfo.UserInfo.Type = "CITIZEN"
	billRequest := map[string]interface{}{
		"RequestInfo": requestInfo,
	}

	// Log request body
	if jsonBytes, err := json.MarshalIndent(billRequest, "", "  "); err == nil {
		log.Printf("FetchBill Request JSON:\n%s", string(jsonBytes))
	} else {
		log.Printf("Failed to marshal billRequest: %v", err)
	}

	// Make HTTP POST call
	var resp demand.BillResponse
	err := r.restCallRepo.Post(url, billRequest, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to call billing service: %w", err)
	}

	return resp.Bill, nil
}
