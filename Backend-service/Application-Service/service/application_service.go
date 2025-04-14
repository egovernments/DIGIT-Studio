package service

import (
	"Application-Service/model"
	"Application-Service/repository"
	"context"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
}

func NewApplicationService(repo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, req model.ApplicationRequest) (model.ApplicationResponse, error) {
	return s.repo.Create(ctx, req)
}

func (s *ApplicationService) SearchApplication(ctx context.Context, criteria model.SearchCriteria) (model.SearchResponse, error) {
	return s.repo.Search(ctx, criteria)
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, req model.ApplicationRequest) (model.ApplicationResponse, error) {
	return s.repo.Update(ctx, req)
}
