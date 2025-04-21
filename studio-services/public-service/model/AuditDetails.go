package model

import (
	"github.com/google/uuid"
	"math/big"
)

type AuditDetails struct {
	CreatedBy        uuid.UUID `json:"createdBy"`
	LastModifiedBy   uuid.UUID `json:"lastModifiedBy"`
	CreatedTime      big.Int   `json:"createdTime"`
	LastModifiedTime big.Int   `json:"lastModifiedTime"`
}
