package model

type SearchCriteria struct {
	TenantId          string      `json:"tenantId"`
	Ids               []string    `json:"ids"`
	BusinessService   string      `json:"businessService"`
	Module            string      `json:"module"`
	ReferenceId       []Reference `json:"referenceId"`
	ApplicationNumber string      `json:"applicationNumber"`
	Status            string      `json:"status"`
	ServiceCode       string      `json:"serviceCode"`
}
