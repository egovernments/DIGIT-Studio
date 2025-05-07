package service

import (
	"log"
	"public-service/config"
	"public-service/model"
	"public-service/repository"
)

type MDMSService struct {
	restCallRepo repository.RestCallRepository
}

func NewMDMSService(repo repository.RestCallRepository) *MDMSService {
	return &MDMSService{
		restCallRepo: repo,
	}
}

func (s *MDMSService) MDMSSearch(criteria model.MdmsCriteria, info model.RequestInfo) (model.MdmsResponse, error) {
	log.Println("MdmsCriteria:", criteria)

	url := config.GetEnv("MDMS_SERVICE_HOST") + config.GetEnv("MDMS_SEARCH_ENDPOINT")
	log.Println("MDMS service URL: " + url)

	// Create the request body
	req := map[string]interface{}{
		"RequestInfo":  info,
		"MdmsCriteria": criteria,
	}

	var resp model.MdmsResponse

	err := s.restCallRepo.Post(url, req, &resp)
	if err != nil {
		log.Printf("Error calling MDMS service: %v", err)
		return model.MdmsResponse{}, err
	}

	log.Println("MDMS Response:", resp)
	return resp, nil
}
