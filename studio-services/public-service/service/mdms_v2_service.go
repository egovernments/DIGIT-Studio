package service

import (
	"log"
	"os"
	"public-service/model"
	"public-service/repository"
)

type MDMSV2Service struct {
	restCallRepo repository.RestCallRepository
}

func NewMDMSV2Service(repo repository.RestCallRepository) *MDMSV2Service {
	return &MDMSV2Service{
		restCallRepo: repo,
	}
}

func (s *MDMSV2Service) SearchMDMS(tenantId, schemaCode, serviceName, module string, requestInfo model.RequestInfo) (map[string]interface{}, error) {

	url := os.Getenv("MDMS_SERVICE_HOST") + os.Getenv("MDMS_V2_SEARCH_ENDPOINT")

	payload := model.MDMSV2Request{
		MdmsCriteria: model.MdmsV2Criteria{
			TenantID:   tenantId,
			Filters:    map[string]string{"service": serviceName, "module": module},
			SchemaCode: schemaCode,
			Limit:      10,
			Offset:     0,
		},
		RequestInfo: requestInfo,
	}
	log.Println("url", url)
	log.Println("payload", payload)
	//reqBody, err := json.Marshal(payload)

	var resp map[string]interface{}
	err := s.restCallRepo.Post(url, payload, &resp)
	if err != nil {
		log.Printf("Error calling MDMS service: %v", err)
		return nil, err
	}

	log.Println("MDMS Response:", resp)
	return resp, nil
}
