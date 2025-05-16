package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

/*
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
		url := fmt.Sprintf("%s%s%s?locale=%s&tenantId=%s&module=digit-studio&codes=%s",
			os.Getenv("LOCALIZATION_SERVICE_HOST"),
			os.Getenv("LOCALIZATION_CONTEXT_PATH"),
			os.Getenv("LOCALIZATION_SEARCH_ENDPOINT"),
			locale,
			tenantID, // ensure tenantID has at least 2 chars
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
*/
func (l *LocalizationService) GetLocalizationMessage(requestInfo model.RequestInfo, code string, tenantID string) map[string]string {
	msgDetail := make(map[string]string)
	locale := os.Getenv("NOTIFICATION_LOCALE")

	if requestInfo.MsgId != "" {
		parts := strings.Split(requestInfo.MsgId, "|")
		if len(parts) >= 2 {
			locale = parts[1]
		}
	}

	url := fmt.Sprintf("%s%s%s?locale=%s&tenantId=%s&module=digit-studio&codes=%s",
		os.Getenv("LOCALIZATION_SERVICE_HOST"),
		os.Getenv("LOCALIZATION_CONTEXT_PATH"),
		os.Getenv("LOCALIZATION_SEARCH_ENDPOINT"),
		locale,
		tenantID,
		code,
	)

	reqBody := map[string]interface{}{
		"RequestInfo": requestInfo,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		return msgDetail
	}

	// 🔍 Print request body as pretty JSON
	var prettyReq bytes.Buffer
	if err := json.Indent(&prettyReq, bodyBytes, "", "  "); err == nil {
		log.Printf("Localization Request Body:\n%s", prettyReq.String())
	} else {
		log.Printf("Failed to format request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Error calling localization service: %v", err)
		return msgDetail
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return msgDetail
	}

	// 🔍 Print response body as pretty JSON
	var prettyResp bytes.Buffer
	if err := json.Indent(&prettyResp, respBytes, "", "  "); err == nil {
		log.Printf("Localization Response Body:\n%s", prettyResp.String())
	} else {
		log.Printf("Failed to format response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		log.Printf("Error decoding localization response: %v", err)
		return msgDetail
	}

	if messages, ok := result["messages"].([]interface{}); ok && len(messages) > 0 {
		firstMsg := messages[0]
		if msgMap, ok := firstMsg.(map[string]interface{}); ok {
			if msgStr, ok := msgMap["message"].(string); ok {
				msgDetail["message"] = msgStr
			} else {
				log.Printf("Message field missing or not a string in first message: %v", firstMsg)
			}
		} else {
			log.Printf("First message is not a valid map: %v", firstMsg)
		}
	} else {
		log.Printf("No messages found in localization response.")
	}

	msgDetail["templateId"] = ""
	return msgDetail
}
