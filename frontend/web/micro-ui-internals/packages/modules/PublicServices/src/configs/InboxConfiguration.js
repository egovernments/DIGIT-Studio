import React from "react";

export const InboxConfig = (module) => {
    module = module.toUpperCase();

    return ({
        headerLabel: "Inbox",
        type: "inbox",
        // apiDetails: {
        //     serviceName: `/egov-mdms-service/v2/_search`,
        //     requestParam: {},
        //     requestBody: {
        //         MdmsCriteria: {},
        //     },
        //     minParametersForSearchForm: 0,
        //     masterName: "commonUiConfig",
        //     moduleName: "SearchMDMSConfig",
        //     tableFormJsonPath: "requestBody.MdmsCriteria",
        //     filterFormJsonPath: "requestBody.MdmsCriteria.custom",
        //     searchFormJsonPath: "requestBody.MdmsCriteria.custom",
        // },
        apiDetails: {
            serviceName: "/inbox-v2/v2/_search",
            requestParam: {},
            requestBody: {
                inbox: {
                    processSearchCriteria: {
                        businessService: [
                            "PGR"
                        ],
                        moduleName: "RAINMAKER-PGR"
                    },
                    moduleSearchCriteria: {}
                }
            },
            minParametersForSearchForm: 0,
            minParametersForFilterForm: 0,
            masterName: "commonUiConfig",
            moduleName: "PGRInboxConfig",
            tableFormJsonPath: "requestBody.inbox",
            filterFormJsonPath: "requestBody.inbox.moduleSearchCriteria",
            searchFormJsonPath: "requestBody.inbox.moduleSearchCriteria",
        },
        // apiDetails: {
        //     serviceName: "/mdms-v2/v2/_search",
        //     requestParam: {},
        //     requestBody: {
        //         inbox: {
        //             "SearchCriteria": {
        //             }
        //         }
        //     },
        //     minParametersForSearchForm: 0,
        //     minParametersForFilterForm: 0,
        //     masterName: "commonUiConfig",
        //     moduleName: "SampleInboxConfig",
        //     tableFormJsonPath: "requestBody.inbox",
        //     filterFormJsonPath: "requestBody.inbox.SearchCriteria",
        //     searchFormJsonPath: "requestBody.inbox.SearchCriteria",
        // },
        sections: {
            search: {
                uiConfig: {
                    headerStyle: {},
                    formClassName: "custom-digit--search-field-wrapper-classname",
                    primaryLabel: "ES_COMMON_SEARCH",
                    secondaryLabel: "ES_COMMON_CLEAR_SEARCH",
                    minReqFields: 1,
                    defaultValues: {
                    },
                    fields: [
                        {
                            inline: true,
                            label: "Application Number",
                            isMandatory: false,
                            type: "text",
                            disable: false,
                            populators: { name: "application", error: "Error!" },
                        },
                        {
                            label: "Business Service",
                            isMandatory: true,
                            key: "service",
                            type: "dropdown",
                            disable: false,
                            preProcess: {
                                updateDependent: ["populators.options"]
                            },
                            populators: {
                                name: "service",
                                optionsKey: "code",
                                options: [
                                ]
                            }
                        },
                    ],
                },
                label: "",
                show: true,
            },
            links: {
                uiConfig: {
                    links: [
                    ],
                    label: "",
                },
                children: {},
                show: true,
            },
            filter: {
                uiConfig: {
                    type: "filter",
                    headerStyle: null,
                    primaryLabel: "Filter",
                    secondaryLabel: "",
                    minReqFields: 1,
                    defaultValues: {
                        state: "",
                        ward: [],
                        locality: [],
                        assignee: {
                            code: "ASSIGNED_TO_ALL",
                            name: "EST_INBOX_ASSIGNED_TO_ALL",
                        },
                    },
                    fields: [
                        {
                            label: "",
                            type: "radio",
                            isMandatory: false,
                            disable: false,
                            populators: {
                                name: "assignee",
                                options: [
                                    {
                                        code: "ASSIGNED_TO_ME",
                                        name: `${module}_ASSIGNED_TO_ME`,
                                    },
                                    {
                                        code: "ASSIGNED_TO_ALL",
                                        name: `${module}_ASSIGNED_TO_ALL`,
                                    },
                                ],
                                optionsKey: "name",
                            },
                        },
                        {
                            label: "COMMON_WARD",
                            type: "locationdropdown",
                            isMandatory: false,
                            disable: false,
                            populators: {
                                name: "ward",
                                type: "ward",
                                optionsKey: "i18nKey",
                                defaultText: "COMMON_SELECT_WARD",
                                selectedText: "COMMON_SELECTED",
                                allowMultiSelect: true
                            }
                        },
                        {
                            label: "COMMON_WORKFLOW_STATES",
                            type: "workflowstatesfilter",
                            labelClassName: "checkbox-status-filter-label",
                            isMandatory: false,
                            disable: false,
                            populators: {
                                name: "state",
                                labelPrefix: "WF_MUSTOR_",
                                businessService: "NEWTL",
                            },
                        },
                    ],
                },
                label: "ES_COMMON_FILTERS",
                show: true,
            },
            searchResult: {
                uiConfig: {
                    columns: [
                        {
                            label: "Application Number",
                            jsonPath: "",
                        },
                        {
                            label: "Business Service",
                            jsonPath: "",
                        },
                        {
                            label: "Status",
                            jsonPath: "",
                        },
                        {
                            label: "SLA",
                            jsonPath: "",
                        }
                    ],
                    tableProps: {
                        showTableDescription: "This is the search table description",
                        showTableTitle: "Search table title",
                        tableClassName: "custom-classname-resultsdatatable"
                    },
                    actionProps: {
                        actions: [
                            {
                                label: "Action1",
                                variation: "secondary",
                                icon: "Edit",
                            },
                            {
                                label: "Action2",
                                variation: "primary",
                                icon: "CheckCircle",
                            },
                        ],
                    },
                    enableColumnSort: true,
                    resultsJsonPath: "mdms",
                    defaultSortAsc: true,
                },
                children: {},
                show: true,
            },
        },
        // additionalSections: {}, 
    })
};
