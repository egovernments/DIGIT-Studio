import React from "react";

export const searchConfig = (module) => {
    console.log(module.toUpperCase());

    return ({
        headerLabel: "Search",
        type: "search",
        apiDetails: {
            serviceName: `/public-service/v1/application/SVC-DEV-TRADELICENSE-NEWTL-04`,
            requestParam: {
            },
            requestBody: {
                MdmsCriteria: {
                },
            },
            // headers:{
            //     "X-Tenant-Id": "dev"
            // },
            minParametersForSearchForm: 0,
            masterName: "commonUiConfig",
            moduleName: "SearchMDMSConfig",
            tableFormJsonPath: "requestBody.MdmsCriteria",
            filterFormJsonPath: "requestBody.MdmsCriteria.filters",
            searchFormJsonPath: "requestBody.MdmsCriteria.filters",
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
                            label: "Application Number",
                            isMandatory: false,
                            type: "text",
                            disable: false,
                            populators: { name: "application", error: "Error!" },
                        },
                        {
                            inline: true,
                            label: "Status",
                            isMandatory: false,
                            type: "text",
                            disable: false,
                            populators: { name: "status", error: "Error!" },
                        },
                        {
                            inline: true,
                            label: "To Date",
                            isMandatory: false,
                            description: "",
                            type: "date",
                            disable: false,
                            populators: { name: "todate", error: "Error!" },
                        },
                        {
                            inline: true,
                            label: "From Date",
                            isMandatory: false,
                            description: "",
                            type: "date",
                            disable: false,
                            populators: { name: "fromdate", error: "Error!" },
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
                                options: [],
                            }
                        },
                    ],
                },
                label: "",
                show: true,
            },
            searchResult: {
                uiConfig: {
                    columns: [
                        {
                            label: "Application",
                            jsonPath: "uniqueIdentifier",
                        },
                        {
                            label: "Status",
                            jsonPath: "",
                        },
                        {
                            label: "To Date",
                            jsonPath: "",
                        },
                        {
                            label: "From Date",
                            jsonPath: "",
                        },
                        {
                            label: "Business Service",
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
