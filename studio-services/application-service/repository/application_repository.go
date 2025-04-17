package repository

import (
	"Application-Service/model"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"math/big"
	"strings"
	"time"
)

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// Create Application
func (r *ApplicationRepository) Create(ctx context.Context, req model.ApplicationRequest) (model.ApplicationResponse, error) {
	searchCriteria := model.SearchCriteria{
		TenantId:        fmt.Sprintf("%d", req.Application.TenantId),
		Module:          req.Application.Module,
		BusinessService: req.Application.BusinessService,
	}

	existingApps, _ := r.Search(ctx, searchCriteria)
	if len(existingApps.Application) > 0 {
		return model.ApplicationResponse{}, errors.New("application already exists")
	}

	now := time.Now()
	createdBy := req.RequestInfo.UserInfo.UserId
	appID := uuid.New()

	// Set missing IDs
	req.Application.Address.Id = uuid.New()
	req.Application.Workflow.Id = uuid.New()
	req.RequestInfo.UserInfo.Id = uuid.New()

	// Always generate new Application Number
	req.Application.ApplicationNumber, _ = r.generateApplicationNumber(ctx, req.Application.TenantId, req.Application.Module, req.Application.BusinessService)

	// Marshal complex fields
	serviceDetailsJSON, _ := json.Marshal(req.Application.ServiceDetails)
	additionalDetailsJSON, _ := json.Marshal(req.Application.AdditionalDetails)
	addressJSON, _ := json.Marshal(req.Application.Address)
	workflowJSON, _ := json.Marshal(req.Application.Workflow)

	insertQuery := `
		INSERT INTO eg_applications (
			id, tenant_id, module, business_service, status, channel, application_number,
			workflow_status, service_details, additional_details, address, workflow,
			createdby, last_modifiedby, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14, $15, $16
		)
	`
	_, err := r.db.ExecContext(ctx, insertQuery,
		appID,
		req.Application.TenantId,
		req.Application.Module,
		req.Application.BusinessService,
		req.Application.Status,
		req.Application.Channel,
		req.Application.ApplicationNumber,
		req.Application.WorkflowStatus,
		serviceDetailsJSON,
		additionalDetailsJSON,
		addressJSON,
		workflowJSON,
		createdBy,
		createdBy,
		now,
		now,
	)
	if err != nil {
		return model.ApplicationResponse{}, err
	}

	// Insert References
	for i, ref := range req.Application.Reference {
		refID := uuid.New()
		refQuery := `
			INSERT INTO eg_reference (
				id, reference_type, module, tenant_id, reference_no, active, application_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err := r.db.ExecContext(ctx, refQuery,
			refID,
			ref.ReferenceType,
			ref.Module,
			ref.TenantId,
			ref.ReferenceNo,
			ref.Active,
			appID,
		)
		if err != nil {
			return model.ApplicationResponse{}, err
		}
		req.Application.Reference[i].Id = refID
	}

	// Insert Applicants
	for i, applicant := range req.Application.Applicants {
		applicantID := uuid.New()
		applicantQuery := `
			INSERT INTO eg_applicant (
				id, type, application_id, user_id, name, mobile_number, email_id, prefix, active
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
		_, err := r.db.ExecContext(ctx, applicantQuery,
			applicantID,
			applicant.Type,
			appID,
			applicant.UserId,
			applicant.Name,
			applicant.MobileNumber,
			applicant.EmailId,
			applicant.Prefix,
			applicant.Active,
		)
		if err != nil {
			return model.ApplicationResponse{}, err
		}
		req.Application.Applicants[i].Id = applicantID
	}

	nowMillis := time.Now().UnixMilli()

	return model.ApplicationResponse{
		ResponseInfo: model.ResponseInfo{
			ApiId: req.RequestInfo.ApiId,
			Ver:   req.RequestInfo.Ver,
		},
		Application: model.Application{
			Id:                appID,
			TenantId:          req.Application.TenantId,
			Module:            req.Application.Module,
			BusinessService:   req.Application.BusinessService,
			Status:            req.Application.Status,
			Channel:           req.Application.Channel,
			ApplicationNumber: req.Application.ApplicationNumber,
			WorkflowStatus:    req.Application.WorkflowStatus,
			ServiceDetails:    req.Application.ServiceDetails,
			AdditionalDetails: req.Application.AdditionalDetails,
			Address:           req.Application.Address,
			Workflow:          req.Application.Workflow,
			Applicants:        req.Application.Applicants,
			Reference:         req.Application.Reference,
			AuditDetails: model.AuditDetails{
				CreatedBy:        createdBy,
				LastModifiedBy:   createdBy,
				CreatedTime:      *big.NewInt(nowMillis),
				LastModifiedTime: *big.NewInt(nowMillis),
			},
		},
	}, nil
}
func (r *ApplicationRepository) Search(ctx context.Context, criteria model.SearchCriteria) (model.SearchResponse, error) {
	var queryBuilder strings.Builder
	var args []interface{}
	var conditions []string
	argPos := 1

	queryBuilder.WriteString(`
		SELECT 
			a.id, a.tenant_id, a.module, a.business_service, a.status, a.channel, a.application_number,
			a.workflow_status, a.service_details, a.additional_details, a.address, a.workflow,
			a.createdby, a.last_modifiedby, a.created_at, a.updated_at,
			r.id, r.reference_type, r.module, r.tenant_id, r.reference_no, r.active,
			ap.id, ap.type, ap.user_id, ap.name, ap.mobile_number, ap.email_id, ap.prefix, ap.active
		FROM eg_applications a
		LEFT JOIN eg_reference r ON a.id = r.application_id
		LEFT JOIN eg_applicant ap ON a.id = ap.application_id
	`)

	// Dynamic where clauses
	if criteria.TenantId != "" {
		conditions = append(conditions, fmt.Sprintf("a.tenant_id = $%d", argPos))
		args = append(args, criteria.TenantId)
		argPos++
	}
	if len(criteria.Ids) > 0 {
		conditions = append(conditions, fmt.Sprintf("a.id = ANY($%d)", argPos))
		args = append(args, pq.Array(criteria.Ids))
		argPos++
	}
	if criteria.Module != "" {
		conditions = append(conditions, fmt.Sprintf("a.module = $%d", argPos))
		args = append(args, criteria.Module)
		argPos++
	}
	if criteria.BusinessService != "" {
		conditions = append(conditions, fmt.Sprintf("a.business_service = $%d", argPos))
		args = append(args, criteria.BusinessService)
		argPos++
	}
	if criteria.ApplicationNo != "" {
		conditions = append(conditions, fmt.Sprintf("a.application_number = $%d", argPos))
		args = append(args, criteria.ApplicationNo)
		argPos++
	}
	if criteria.Status != "" {
		conditions = append(conditions, fmt.Sprintf("a.status = $%d", argPos))
		args = append(args, criteria.Status)
		argPos++
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	log.Println("query:", queryBuilder.String())
	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return model.SearchResponse{}, err
	}
	defer rows.Close()

	var applications []model.Application
	appMap := make(map[uuid.UUID]*model.Application)

	for rows.Next() {
		var (
			appId                                                                                 uuid.UUID
			tenantId, module, businessService, status, channel, applicationNumber, workflowStatus string
			serviceDetailsJSON, additionalDetailsJSON, addressJSON, workflowJSON                  []byte
			createdBy, lastModifiedBy                                                             string
			createdAt, updatedAt                                                                  time.Time
			refId                                                                                 sql.NullString
			ref                                                                                   model.Reference
			applicantId                                                                           sql.NullString
			applicant                                                                             model.User
		)

		err := rows.Scan(
			&appId,
			&tenantId,
			&module,
			&businessService,
			&status,
			&channel,
			&applicationNumber,
			&workflowStatus,
			&serviceDetailsJSON,
			&additionalDetailsJSON,
			&addressJSON,
			&workflowJSON,
			&createdBy,
			&lastModifiedBy,
			&createdAt,
			&updatedAt,
			&refId,
			&ref.ReferenceType,
			&ref.Module,
			&ref.TenantId,
			&ref.ReferenceNo,
			&ref.Active,
			&applicantId,
			&applicant.Type,
			&applicant.UserId,
			&applicant.Name,
			&applicant.MobileNumber,
			&applicant.EmailId,
			&applicant.Prefix,
			&applicant.Active,
		)
		if err != nil {
			return model.SearchResponse{}, err
		}

		app, exists := appMap[appId]
		if !exists {
			app = &model.Application{
				Id:                appId,
				TenantId:          tenantId,
				Module:            module,
				BusinessService:   businessService,
				Status:            status,
				Channel:           channel,
				ApplicationNumber: applicationNumber,
				WorkflowStatus:    workflowStatus,
				AuditDetails: model.AuditDetails{
					CreatedBy:        createdBy,
					LastModifiedBy:   lastModifiedBy,
					CreatedTime:      *big.NewInt(createdAt.UnixMilli()),
					LastModifiedTime: *big.NewInt(updatedAt.UnixMilli()),
				},
			}

			// Unmarshal JSON fields
			_ = json.Unmarshal(serviceDetailsJSON, &app.ServiceDetails)
			_ = json.Unmarshal(additionalDetailsJSON, &app.AdditionalDetails)
			_ = json.Unmarshal(addressJSON, &app.Address)
			_ = json.Unmarshal(workflowJSON, &app.Workflow)

			appMap[appId] = app
		}

		// Add Reference if present
		if refId.Valid {
			ref.Id, _ = uuid.Parse(refId.String)
			app.Reference = append(app.Reference, ref)
		}

		// Add Applicant if present
		if applicantId.Valid {
			applicant.Id, _ = uuid.Parse(applicantId.String)
			app.Applicants = append(app.Applicants, applicant)
		}
	}

	for _, app := range appMap {
		applications = append(applications, *app)
	}

	return model.SearchResponse{
		Application: applications,
		ResponseInfo: model.ResponseInfo{
			Status: "successful",
		},
	}, nil
}

func (r *ApplicationRepository) Update(ctx context.Context, req model.ApplicationRequest) (model.ApplicationResponse, error) {
	nowMillis := time.Now().UnixMilli()
	log.Println("Application request", req.Application)
	// Marshal complex fields
	serviceDetailsJSON, _ := json.Marshal(req.Application.ServiceDetails)
	additionalDetailsJSON, _ := json.Marshal(req.Application.AdditionalDetails)
	addressJSON, _ := json.Marshal(req.Application.Address)
	workflowJSON, _ := json.Marshal(req.Application.Workflow)

	appQuery := `
		UPDATE eg_applications
		SET tenant_id = $1,
		    module = $2,
		    business_service = $3,
		    status = $4,
		    channel = $5,
		    application_number = $6,
		    workflow_status = $7,
		    service_details = $8,
		    additional_details = $9,
		    address = $10,
		    workflow = $11,
		    last_modifiedby = $12,
		    updated_at = to_timestamp($13 / 1000.0)
		WHERE id = $14
	`
	_, err := r.db.ExecContext(ctx, appQuery,
		req.Application.TenantId,
		req.Application.Module,
		req.Application.BusinessService,
		req.Application.Status,
		req.Application.Channel,
		req.Application.ApplicationNumber,
		req.Application.WorkflowStatus,
		serviceDetailsJSON,
		additionalDetailsJSON,
		addressJSON,
		workflowJSON,
		req.RequestInfo.UserInfo.UserId,
		nowMillis,
		req.Application.Id,
	)
	if err != nil {
		return model.ApplicationResponse{}, fmt.Errorf("failed to update application: %w", err)
	}

	// Update References
	for _, ref := range req.Application.Reference {
		refQuery := `
			UPDATE eg_reference
			SET reference_type = $1,
			    module = $2,
			    tenant_id = $3,
			    reference_no = $4,
			    active = $5
			WHERE id = $6
		`
		_, err := r.db.ExecContext(ctx, refQuery,
			ref.ReferenceType,
			ref.Module,
			ref.TenantId,
			ref.ReferenceNo,
			ref.Active,
			ref.Id,
		)
		if err != nil {
			return model.ApplicationResponse{}, fmt.Errorf("failed to update reference: %w", err)
		}
	}

	// Update Applicants
	for _, applicant := range req.Application.Applicants {
		applicantQuery := `
			UPDATE eg_applicant
			SET type = $1,
			    user_id = $2,
			    name = $3,
			    mobile_number = $4,
			    email_id = $5,
			    prefix = $6,
			    active = $7
			WHERE id = $8
		`
		_, err := r.db.ExecContext(ctx, applicantQuery,
			applicant.Type,
			applicant.UserId,
			applicant.Name,
			applicant.MobileNumber,
			applicant.EmailId,
			applicant.Prefix,
			applicant.Active,
			applicant.Id,
		)
		if err != nil {
			return model.ApplicationResponse{}, fmt.Errorf("failed to update applicant: %w", err)
		}
	}

	return model.ApplicationResponse{
		ResponseInfo: model.ResponseInfo{
			ApiId: req.RequestInfo.ApiId,
			Ver:   req.RequestInfo.Ver,
		},
		Application: req.Application,
	}, nil
}

func (r *ApplicationRepository) generateApplicationNumber(ctx context.Context, tenantId, module, businessService string) (applicationNumber string, err error) {
	var lastNumber int
	tx, txErr := r.db.BeginTx(ctx, nil) // Begin a transaction
	if txErr != nil {
		err = fmt.Errorf("failed to begin transaction: %w", txErr)
		return
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Check if entry exists
	query := `
		SELECT last_number
		FROM eg_application_sequence
		WHERE tenant_id = $1 AND module = $2 AND business_service = $3
		FOR UPDATE
	`
	scanErr := tx.QueryRowContext(ctx, query, tenantId, module, businessService).Scan(&lastNumber)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			insertQuery := `
				INSERT INTO eg_application_sequence (tenant_id, module, business_service, last_number)
				VALUES ($1, $2, $3, 1)
			`
			_, execErr := tx.ExecContext(ctx, insertQuery, tenantId, module, businessService)
			if execErr != nil {
				err = fmt.Errorf("failed to insert sequence: %w", execErr)
				return
			}
			lastNumber = 1
		} else {
			err = fmt.Errorf("failed to get sequence: %w", scanErr)
			return
		}
	} else {
		// If exists, increment last_number
		lastNumber += 1
		updateQuery := `
			UPDATE eg_application_sequence
			SET last_number = $1
			WHERE tenant_id = $2 AND module = $3 AND business_service = $4
		`
		_, updateErr := tx.ExecContext(ctx, updateQuery, lastNumber, tenantId, module, businessService)
		if updateErr != nil {
			err = fmt.Errorf("failed to update sequence: %w", updateErr)
			return
		}
	}

	applicationNumber = fmt.Sprintf("%s-%s-%s-%03d", tenantId, module, businessService, lastNumber)
	return
}
