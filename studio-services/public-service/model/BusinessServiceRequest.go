package model

type BusinessServiceRequest struct {
	RequestInfo     RequestInfo     `json:"RequestInfo"`
	BusinessService BusinessService `json:"BusinessService"`
}
