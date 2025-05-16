// repository/repository.go

package repository

import (
	"database/sql"
	"egov-generic-service/model"

	"encoding/json"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
)

func CreateRecord(req model.GenericInitRequest) (model.GenericInitResponse, error) {
	id := uuid.New()

	additionalDetails, err := json.Marshal(req.Service.AdditionalDetails)
	if err != nil {
		return model.GenericInitResponse{}, err
	}

	query := `
		INSERT INTO egov_service_request (
			id, tenantId, businessService, module, status, additionalDetail, createdTime, lastModifiedTime, createdby, lastmodifiedby
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7, $8)
	`

	_, err = GetDB().Exec(
		query,
		id,
		req.Service.TenantId,
		req.Service.BusinessService,
		req.Service.Module,
		req.Service.Status,
		additionalDetails,
		1, // createdBy
		1, // lastModifiedBy
	)
	if err != nil {
		return model.GenericInitResponse{}, err
	}

	service := model.ServiceResp{
		Id:                id.String(),
		TenantId:          req.Service.TenantId,
		BusinessService:   req.Service.BusinessService,
		Module:            req.Service.Module,
		Status:            model.Status(req.Service.Status),
		AdditionalDetails: req.Service.AdditionalDetails,
		AuditDetails: model.AuditDetails{
			CreatedBy:        "1",
			LastModifiedBy:   "1",
			CreatedTime:      int(time.Now().Unix()),
			LastModifiedTime: int(time.Now().Unix()),
		},
	}

	return model.GenericInitResponse{Services: []model.ServiceResp{service}}, nil
}

func UpdateRecord(req model.GenericInitRequest) (model.GenericInitResponse, error) {
	additionalDetails, err := json.Marshal(req.Service.AdditionalDetails)
	if err != nil {
		return model.GenericInitResponse{}, err
	}

	query := `
		UPDATE egov_service_request
		SET status = $1, additionalDetail = $2, lastModifiedTime = NOW(), lastmodifiedby = $3
		WHERE tenantId = $4 AND businessService = $5 AND module = $6
	`

	_, err = GetDB().Exec(
		query,
		req.Service.Status,
		additionalDetails,
		1, // lastModifiedBy
		req.Service.TenantId,
		req.Service.BusinessService,
		req.Service.Module,
	)
	if err != nil {
		return model.GenericInitResponse{}, err
	}

	service := model.ServiceResp{
		TenantId:          req.Service.TenantId,
		BusinessService:   req.Service.BusinessService,
		Module:            req.Service.Module,
		Status:            model.Status(req.Service.Status),
		AdditionalDetails: req.Service.AdditionalDetails,
		AuditDetails: model.AuditDetails{
			LastModifiedBy:   "1",
			LastModifiedTime: int(time.Now().Unix()),
		},
	}

	return model.GenericInitResponse{Services: []model.ServiceResp{service}}, nil
}

func SearchRecords(req model.GenericInitRequest) ([]model.GenericInitResponse, error) {
	query := `SELECT id, tenantid, businessservice, module, status, additionaldetail, createdtime, lastmodifiedtime, createdby, lastmodifiedby
		FROM egov_service_request
		WHERE tenantid = $1 AND businessservice = $2 AND module = $3
	`
	log.Println("Inside Search: ", query)

	rows, err := GetDB().Query(query, req.Service.TenantId, req.Service.BusinessService, req.Service.Module)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	log.Println("search")
	var results []model.GenericInitResponse

	for rows.Next() {
		var (
			id                                uuid.UUID
			tenantId, businessService, module string
			status                            string
			additionalDetail                  sql.NullString
			createdTime, lastModifiedTime     time.Time
			createdBy, lastModifiedBy         int64
		)

		if err := rows.Scan(
			&id, &tenantId, &businessService, &module, &status, &additionalDetail,
			&createdTime, &lastModifiedTime, &createdBy, &lastModifiedBy,
		); err != nil {
			return nil, err
		}

		var additionalDetailsMap map[string]interface{}
		if additionalDetail.Valid {
			if err := json.Unmarshal([]byte(additionalDetail.String), &additionalDetailsMap); err != nil {
				log.Printf("Error unmarshalling additional details: %v", err)
			}
		}

		service := model.ServiceResp{
			Id:                id.String(),
			TenantId:          tenantId,
			BusinessService:   businessService,
			Module:            module,
			Status:            model.Status(status),
			AdditionalDetails: additionalDetailsMap,
			AuditDetails: model.AuditDetails{
				CreatedBy:        strconv.FormatInt(createdBy, 10),
				LastModifiedBy:   strconv.FormatInt(lastModifiedBy, 10),
				CreatedTime:      int(createdTime.Unix()),
				LastModifiedTime: int(lastModifiedTime.Unix()),
			},
		}

		res := model.GenericInitResponse{Services: []model.ServiceResp{service}}
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
