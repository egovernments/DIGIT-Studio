package model

import "math/big"

type AuditDetails struct {
	CreatedBy        string  `json:"createdBy"`
	LastModifiedBy   string  `json:"lastModifiedBy"`
	CreatedTime      big.Int `json:"createdTime"`
	LastModifiedTime big.Int `json:"lastModifiedTime"`
}
