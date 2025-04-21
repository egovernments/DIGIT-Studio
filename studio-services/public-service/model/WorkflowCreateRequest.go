package model

type WorkflowCreateRequest struct {
	RequestInfo      RequestInfo       `json:"RequestInfo"`
	BusinessServices []BusinessService `json:"BusinessServices"`
}

type BusinessService struct {
	TenantId           string  `json:"tenantId"`
	BusinessService    string  `json:"businessService"`
	Business           string  `json:"business"`
	BusinessServiceSla int64   `json:"businessServiceSla"`
	States             []State `json:"states"`
}

type State struct {
	Sla               *int64   `json:"sla"`
	State             *string  `json:"state"`
	ApplicationStatus string   `json:"applicationStatus"`
	DocUploadRequired bool     `json:"docUploadRequired"`
	IsStartState      bool     `json:"isStartState"`
	IsTerminateState  bool     `json:"isTerminateState"`
	IsStateUpdatable  bool     `json:"isStateUpdatable"`
	Actions           []Action `json:"actions"`
}

type Action struct {
	Action    string   `json:"action"`
	NextState string   `json:"nextState"`
	Roles     []string `json:"roles"`
}

func BuildWorkflowRequestFromService(req ServiceRequest) (*WorkflowCreateRequest, bool) {
	_, ok := req.Service.AdditionalDetails["workflow"]
	if !ok {
		return nil, false
	}

	businessService := BusinessService{
		TenantId:           req.Service.TenantId,
		BusinessService:    req.Service.BusinessService,
		Business:           req.Service.Module,
		BusinessServiceSla: 259200000, // You can make this configurable if needed
		States: []State{
			{
				ApplicationStatus: "ACTIVE",
				DocUploadRequired: false,
				IsStartState:      true,
				IsTerminateState:  false,
				IsStateUpdatable:  true,
				Actions: []Action{
					{
						Action:    "APPLIED",
						NextState: "APPROVED",
						Roles:     []string{"CITIZEN", "EMPLOYEE"},
					},
				},
			},
		},
	}

	wfRequest := &WorkflowCreateRequest{
		RequestInfo:      req.RequestInfo,
		BusinessServices: []BusinessService{businessService},
	}

	return wfRequest, true
}
