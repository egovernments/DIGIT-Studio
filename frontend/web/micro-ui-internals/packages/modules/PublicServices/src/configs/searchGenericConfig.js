
export const searchConfig = {
    headerLabel: "Search",
    type: "search",
    apiDetails: {
            serviceName: `/public-service/v1/application`,
            requestParam: {},
            requestBody: {
                SearchCriteria: {},
            },
            minParametersForSearchForm: 0,
            masterName: "commonUiConfig",
            moduleName: "searchGenericConfig",
            tableFormJsonPath: "requestBody.SearchCriteria",
            filterFormJsonPath: "requestBody.SearchCriteria.custom",
            searchFormJsonPath: "requestBody.SearchCriteria.custom",
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
                        populators: { name: "applicationNumber", error: "Error!" },
                    },
                    {
                        inline: true,
                        label: "Status",
                        isMandatory: false,
                        type: "text",
                        disable: true,
                        populators: { name: "status", error: "Error!" },
                    },
                    {
                        inline: true,
                        label: "To Date",
                        isMandatory: false,
                        description: "",
                        type: "date",
                        disable: true,
                        populators: { name: "todate", error: "Error!" },
                    },
                    {
                        inline: true,
                        label: "From Date",
                        isMandatory: false,
                        description: "",
                        type: "date",
                        disable: true,
                        populators: { name: "fromdate", error: "Error!" },
                    },
                    {
                        label: "Business Service",
                        isMandatory: true,
                        key: "businessService",
                        type: "dropdown",
                        disable: false,
                        preProcess : {
                            updateDependent : ["populators.options"]
                        },
                        populators: {
                            name: "businessService",
                            optionsKey:"name",
                            options:[
                            ]
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
                        label: "Application Number",
                        jsonPath: "applicationNumber",
                        additionalCustomization: true,
                    },
                    {
                        label: "Status",
                        jsonPath: "workflowStatus",
                    },
                    {
                        label: "businessService",
                        jsonPath: "businessService",
                    },
                    {
                        label: "serviceCode",
                        jsonPath: "serviceCode", 
                    }
                ],
                tableProps: {
                    //showTableDescription: "This is the search table description",
                    //showTableTitle: "Search table title",
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
                resultsJsonPath: "Application",
                defaultSortAsc: true,
            },
            children: {},
            show: true,
        },
    },
    // additionalSections: {}, 
};
