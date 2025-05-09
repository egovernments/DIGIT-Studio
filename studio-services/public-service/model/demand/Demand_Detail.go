package demand

import (
	"encoding/json"
	"math/big"
	"public-service/model"
)

type DemandDetail struct {
	ID                string              `json:"id,omitempty"`
	DemandID          string              `json:"demandId,omitempty"`
	TaxHeadMasterCode string              `json:"taxHeadMasterCode"` // required
	TaxAmount         *big.Float          `json:"taxAmount"`         // required
	CollectionAmount  *big.Float          `json:"collectionAmount"`  // required
	AuditDetails      *model.AuditDetails `json:"auditDetails,omitempty"`
	TenantID          string              `json:"tenantId,omitempty"`
}

// Custom unmarshaler to set default value for CollectionAmount if not provided
func (d *DemandDetail) UnmarshalJSON(data []byte) error {
	type Alias DemandDetail
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if d.CollectionAmount == nil {
		d.CollectionAmount = big.NewFloat(0)
	}

	return nil
}
