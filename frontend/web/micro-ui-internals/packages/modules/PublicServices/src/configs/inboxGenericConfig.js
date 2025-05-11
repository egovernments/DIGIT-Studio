import React from "react";
import { useParams } from "react-router-dom/cjs/react-router-dom";

export const InboxConfig = () => {
    const { module } = useParams();
    const prefix = `${module?.toUpperCase()}`;

    return ({
        headerLabel: `${prefix}_INBOX_HEADER`,
        type: "inbox",
        apiDetails: {
            serviceName: "/inbox/v2/_search",
            requestParam: {},
            requestBody: {
                inbox: {
                    processSearchCriteria: {
                        businessService: [],
                        moduleName: "public-service",
                    },
                    moduleSearchCriteria: {
                        sortOrder: "ASC",
                        module: "public-service",
                    }
                }
            },
            minParametersForSearchForm: 0,
            minParametersForFilterForm: 0,
            masterName: "commonUiConfig",
            moduleName: "InboxGenericConfig",
            tableFormJsonPath: "requestBody.inbox",
            filterFormJsonPath: "requestBody.inbox.moduleSearchCriteria",
            searchFormJsonPath: "requestBody.inbox.moduleSearchCriteria",
        },
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
                            label: `${prefix}_APPLICATION_NUMBER`,
                            isMandatory: false,
                            type: "text",
                            disable: false,
                            populators: { name: "applicationNumber", error: "Error!" },
                        },
                        {
                            label: `${prefix}_BUSINESS_SERVICE`,
                            isMandatory: true,
                            key: "businessService",
                            type: "dropdown",
                            disable: false,
                            preProcess: {
                                updateDependent: ["populators.options"]
                            },
                            populators: {
                                name: "businessService",
                                optionsKey: "name",
                                options: []
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
                    primaryLabel: `${prefix}_FILTER`,
                    secondaryLabel: "",
                    minReqFields: 1,
                    defaultValues: {
                        state: "",
                        ward: [],
                        locality: [],
                        assignee: {
                            code: `${prefix}_ASSIGNED_TO_ALL`,
                            name: `${prefix}_ASSIGNED_TO_ALL`,
                        }
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
                                        code: `${prefix}_ASSIGNED_TO_ME`,
                                        name: `${prefix}_ASSIGNED_TO_ME`,
                                    },
                                    {
                                        code: `${prefix}_ASSIGNED_TO_ALL`,
                                        name: `${prefix}_ASSIGNED_TO_ALL`,
                                    },
                                ],
                                optionsKey: "name",
                            },
                        },
                        // {
                        //     key: "boundaryComponent",
                        //     type: "boundary",
                        //     inline: true,
                        //     label: "BoundaryFIlter",
                        //     disable: false,
                        //     populators: {
                        //         name: "boundaryComponent",
                        //         levelConfig: { lowestLevel: "LOCALITY", highestLevel: "PROVINCE", isSingleSelect: [] },
                        //         hierarchyType: "NEWTEST00222",
                        //         noCardStyle: true,
                        //         preSelected: ["NEWTEST00222_MO", "NEWTEST00222_MO_11_MARYLAND", "NEWTEST00222_MO_11_06_PLEEBO"],                    // "frozenData":
                        //         // [{
                        //         //     code: "NEWTEST00222_MO",
                        //         //     name: "NEWTEST00222_MO"
                        //         //   },
                        //         //   {
                        //         //     code: "NEWTEST00222_MO.NEWTEST00222_MO_11_MARYLAND",
                        //         //     name: "NEWTEST00222_MO_11_MARYLAND"
                        //         //   },
                        //         //   {
                        //         //     code: "NEWTEST00222_MO.NEWTEST00222_MO_11_MARYLAND.NEWTEST00222_MO_11_06_PLEEBO",
                        //         //     name: "NEWTEST00222_MO_11_06_PLEEBO"
                        //         //   }]
                        //     },
                        // },
                        {
                            label: `${prefix}_COMMON_WARD`,
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
                            label: `${prefix}_COMMON_WORKFLOW_STATES`,
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
                label: `ES_COMMON_FILTERS`,
                show: true,
            },
            searchResult: {
                uiConfig: {
                    columns: [
                        {
                            label: "Application Number",
                            jsonPath: "applicationNumber",
                        },
                        {
                            label: "Business Service",
                            jsonPath: "businessService",
                        },
                        {
                            label: "Status",
                            jsonPath: "statusMap",
                        },
                        {
                            label: "SLA",
                            jsonPath: "nearingSlaCount",
                        }
                    ],
                    tableProps: {
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
                    resultsJsonPath: "items",
                    defaultSortAsc: true,
                },
                children: {},
                show: true,
            },
        },
    })
};