package service

import (
	"egov-generic-service/model"
	"egov-generic-service/repository"
	"errors"
	"log"
)

func InitGenericService(req model.GenericInitRequest) (model.GenericInitResponse, error) {
	// First search if the record already exists
	log.Println("Inside...")
	existingRecords, err := repository.SearchRecords(req)
	if err != nil {
		return model.GenericInitResponse{}, err
	}

	if len(existingRecords) > 0 {
		// Record exists, return an error or return the existing record
		return model.GenericInitResponse{}, errors.New("record already exists with given tenantId, businessService, and module")
	}

	// If no record found, create new record
	return repository.CreateRecord(req)
}

func UpdateGenericService(req model.GenericInitRequest) (model.GenericInitResponse, error) {
	return repository.UpdateRecord(req)
}

func SearchGenericService(req model.GenericInitRequest) ([]model.GenericInitResponse, error) {
	return repository.SearchRecords(req)
}
