package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/oliveagle/jsonpath"
	"log"
	"net/http"
	"os"
	"public-service/model"
	"public-service/repository"
	"strings"
)

type LocalizationService struct {
	restCallRepo repository.RestCallRepository
}

func NewLocalizationService(repo repository.RestCallRepository) *LocalizationService {
	return &LocalizationService{
		restCallRepo: repo,
	}
}

func (l *LocalizationService) GetLocalizationMessage(requestInfo model.RequestInfo, code string, tenantID string) map[string]string {
	msgDetail := make(map[string]string)
	locale := os.Getenv("NOTIFICATION_LOCALE")

	if requestInfo.MsgId != "" {
		parts := strings.Split(requestInfo.MsgId, "|")
		if len(parts) >= 2 {
			locale = parts[1]
		}
	}

	// Build URL
	url := fmt.Sprintf("%s%s%s?locale=%s&tenantId=%s&module=mgramseva-common&codes=%s",
		os.Getenv("LOCALIZATION_SERVICE_HOST"),
		os.Getenv("LOCALIZATION_CONTEXT_PATH"),
		os.Getenv("LOCALIZATION_SEARCH_ENDPOINT"),
		locale,
		tenantID[:2], // ensure tenantID has at least 2 chars
		code,
	)

	// Create request body
	reqBody := map[string]interface{}{
		"RequestInfo": requestInfo,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return msgDetail
	}

	// Send POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Error calling localization service: %v", err)
		return msgDetail
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding localization response: %v", err)
		return msgDetail
	}

	// Extract message using oliveagle/jsonpath
	message := ""
	jsonPathExpr, err := jsonpath.Compile("$.messages[0].message")
	if err != nil {
		log.Printf("Error compiling jsonpath: %v", err)
	} else {
		if val, err := jsonPathExpr.Lookup(result); err == nil {
			if msgStr, ok := val.(string); ok {
				message = msgStr
			}
		} else {
			log.Printf("Error looking up jsonpath: %v", err)
		}
	}

	msgDetail["message"] = message
	msgDetail["templateId"] = "" // templateId is null in original logic

	return msgDetail
}
