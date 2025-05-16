package model

type GenericInitRequest struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
	Service     Service     `json:"Service"`
}

type RequestInfo struct {
	ApiId       string `json:"apiId"`
	Ver         string `json:"ver"`
	Ts          int    `json:"ts"`
	Action      string `json:"action"`
	Did         string `json:"did"`
	Key         string `json:"key"`
	MsgId       string `json:"msgId"`
	RequesterId string `json:"requesterId"`
	AuthToken   string `json:"authToken"`
}

type Service struct {
	TenantId          string                 `json:"tenantId"`
	BusinessService   string                 `json:"businessService"`
	Module            string                 `json:"module"`
	Status            string                 `json:"status"`
	AdditionalDetails map[string]interface{} `json:"additionalDetails"` // Fixed as dynamic
}
