package model

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

type GenericInitResponse struct {
	ResponseInfo ResponseInfo  `json:"ResponseInfo"`
	Services     []ServiceResp `json:"Services"`
	Pagination   Pagination    `json:"pagination"`
}

type ResponseInfo struct {
	ApiId    string `json:"apiId"`
	Ver      string `json:"ver"`
	Ts       int    `json:"ts"`
	ResMsgId string `json:"resMsgId"`
	MsgId    string `json:"msgId"`
	Status   string `json:"status"`
}

type ServiceResp struct {
	Id                string                 `json:"id"`
	TenantId          string                 `json:"tenantId"`
	BusinessService   string                 `json:"businessService"`
	Module            string                 `json:"module"`
	Status            Status                 `json:"status"`
	AdditionalDetails map[string]interface{} `json:"additionalDetails"` // Fixed as dynamic
	AuditDetails      AuditDetails           `json:"auditDetails"`
}

type AuditDetails struct {
	CreatedBy        string `json:"createdBy"`
	LastModifiedBy   string `json:"lastModifiedBy"`
	CreatedTime      int    `json:"createdTime"`
	LastModifiedTime int    `json:"lastModifiedTime"`
}

type Pagination struct {
	Limit      int    `json:"limit"`
	OffSet     int    `json:"offSet"`
	TotalCount int    `json:"totalCount"`
	SortBy     string `json:"sortBy"`
	Order      string `json:"order"` // Fixed to string ("ASC" / "DESC")
}
