package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"public-service/model"

	"github.com/joho/godotenv"
)

type WorkflowService struct {
	client *http.Client
}

func NewWorkflowService() *WorkflowService {
	return &WorkflowService{
		client: &http.Client{},
	}
}

func (ws *WorkflowService) CreateBusinessService(ctx context.Context, payload model.WorkflowCreateRequest) (*http.Response, error) {
	_ = godotenv.Load() // Load from .env

	host := os.Getenv("WORKFLOW_HOST")
	path := os.Getenv("WORKFLOW_CREATE_PATH")
	url := fmt.Sprintf("%s%s", host, path)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return ws.client.Do(req)
}
