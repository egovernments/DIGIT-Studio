package service

import (
	"context"
	"public-service/model"
	"public-service/repository"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
}

func NewApplicationService(repo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, req model.ApplicationRequest, ServiceCode string) (model.ApplicationResponse, error) {
	return s.repo.Create(ctx, req, ServiceCode)
}

func (s *ApplicationService) SearchApplication(ctx context.Context, criteria model.SearchCriteria) (model.SearchResponse, error) {
	return s.repo.Search(ctx, criteria)
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, req model.ApplicationRequest, serviceCode string, applicationId string) (model.ApplicationResponse, error) {
	return s.repo.Update(ctx, req, serviceCode, applicationId)
}
