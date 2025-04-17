package model

type SearchCriteria struct {
	TenantId        string      `json:"tenantId"`
	Ids             []string    `json:"ids"`
	BusinessService string      `json:"businessService"`
	Module          string      `json:"module"`
	ReferenceId     []Reference `json:"referenceId"`
	ApplicationNo   string      `json:"applicationNo"`
	Status          string      `json:"status"`
}
