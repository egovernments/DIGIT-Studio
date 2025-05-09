import React from "react";
//import { useTranslation } from "react-i18next";

const response = {
    "ResponseInfo": null,
    "ServiceDefinitions": [
        {
            "id": "90a4fc55-c90e-404b-a45f-bfb3a03a473b",
            "tenantId": "mz",
            "code": "IRS_configure-1.WAREHOUSE.DISTRIBUTOR",
            "isActive": true,
            "attributes": [
                {
                    "id": "13e787d3-9fa1-42ab-84fa-c3cd92d0b93d",
                    "referenceId": "90a4fc55-c90e-404b-a45f-bfb3a03a473b",
                    "tenantId": "mz",
                    "code": "SN1",
                    "dataType": "SingleValueList",
                    "values": [
                        "SHORTAGES",
                        "QUALITY_COMPLAINTS",
                        "NOT_SELECTED"
                    ],
                    "isActive": true,
                    "required": false,
                    "regex": null,
                    "order": "1",
                    "auditDetails": {
                        "createdBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "lastModifiedBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "createdTime": 1746526115203,
                        "lastModifiedTime": 1746526115203
                    },
                    "additionalFields": {
                        "fields": [
                            {
                                "key": "007bc6e9-1873-47d7-bc0d-cae51f704544",
                                "value": {
                                    "id": "2d4a7b1e-1f2f-4a8a-9672-43396c6c9a1c",
                                    "key": 1,
                                    "type": {
                                        "code": "SingleValueList"
                                    },
                                    "level": 1,
                                    "title": "Is there a feedback system for health facilities to report any issues or requests related to bednet distribution?",
                                    "value": null,
                                    "options": [
                                        {
                                            "id": "0cff9846-03a2-4453-bf0e-200cdda5f390",
                                            "key": 1,
                                            "label": "Shortages",
                                            "subQuestions": [],
                                            "optionComment": true,
                                            "optionDependency": false,
                                            "parentQuestionId": "2d4a7b1e-1f2f-4a8a-9672-43396c6c9a1c"
                                        },
                                        {
                                            "id": "2d4a7b1e-7c0d-48b1-9d53-8601c6264b90",
                                            "key": 2,
                                            "label": "Quality complaints",
                                            "subQuestions": [],
                                            "optionDependency": false,
                                            "parentQuestionId": "2d4a7b1e-1f2f-4a8a-9672-43396c6c9a1c"
                                        }
                                    ],
                                    "isActive": true,
                                    "parentId": null,
                                    "isRequired": false,
                                    "subQuestions": []
                                }
                            }
                        ],
                        "schema": "serviceDefinition",
                        "version": 1
                    }
                },
                {
                    "id": "c20ddd36-6631-4f97-9832-971b02069dde",
                    "referenceId": "90a4fc55-c90e-404b-a45f-bfb3a03a473b",
                    "tenantId": "mz",
                    "code": "SN2",
                    "dataType": "SingleValueList",
                    "values": [
                        "HOSPITALS",
                        "CLINICS",
                        "COMMUNITY_HEALTH_CENTERS",
                        "NOT_SELECTED"
                    ],
                    "isActive": true,
                    "required": false,
                    "regex": null,
                    "order": "2",
                    "auditDetails": {
                        "createdBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "lastModifiedBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "createdTime": 1746526115203,
                        "lastModifiedTime": 1746526115203
                    },
                    "additionalFields": {
                        "fields": [
                            {
                                "key": "4b43773c-8fd3-4955-9292-62615af6b4bd",
                                "value": {
                                    "id": "4add5323-fc98-4e71-a783-27dbb922c99f",
                                    "key": 2,
                                    "type": {
                                        "code": "SingleValueList"
                                    },
                                    "level": 1,
                                    "title": "What types of health facilities do you distribute to?",
                                    "value": null,
                                    "options": [
                                        {
                                            "id": "34eac43a-e0b5-428f-9d11-12fc5b10b1ac1",
                                            "key": 1,
                                            "label": "Hospitals",
                                            "subQuestions": [],
                                            "optionComment": false,
                                            "optionDependency": false,
                                            "parentQuestionId": "4add5323-fc98-4e71-a783-27dbb922c99f"
                                        },
                                        {
                                            "id": "23ace43b-e0b5-428f-9d11-12fc5b10b1ac1",
                                            "key": 2,
                                            "label": "Clinics",
                                            "subQuestions": [
                                                {
                                                    "id": "c65ac34b-7cc0-4993-a8fe-37e854d2b189",
                                                    "key": 4,
                                                    "type": {
                                                        "code": "SingleValueList"
                                                    },
                                                    "level": 2,
                                                    "title": "Do you have enough products for distribution to health facilities?",
                                                    "value": null,
                                                    "options": [
                                                        {
                                                            "id": "cb45ca84-7cc0-4993-a8fe-37e854d2b189",
                                                            "key": 1,
                                                            "label": "Yes",
                                                            "subQuestions": [],
                                                            "optionDependency": false,
                                                            "parentQuestionId": "c65ac34b-7cc0-4993-a8fe-37e854d2b189"
                                                        },
                                                        {
                                                            "id": "a54c73cb-60da-4c51-8501-cf4a4f473a66",
                                                            "key": 2,
                                                            "label": "No",
                                                            "subQuestions": [],
                                                            "optionComment": true,
                                                            "optionDependency": false,
                                                            "parentQuestionId": "c65ac34b-7cc0-4993-a8fe-37e854d2b189"
                                                        }
                                                    ],
                                                    "isActive": true,
                                                    "parentId": "23ace43b-e0b5-428f-9d11-12fc5b10b1ac1",
                                                    "isRequired": false,
                                                    "subQuestions": []
                                                }
                                            ],
                                            "optionComment": false,
                                            "optionDependency": true,
                                            "parentQuestionId": "4add5323-fc98-4e71-a783-27dbb922c99f"
                                        },
                                        {
                                            "id": "32bbca43-db87-469b-8be4-22012cc22284",
                                            "key": 3,
                                            "label": "Community health centers",
                                            "subQuestions": [],
                                            "optionDependency": false,
                                            "parentQuestionId": "4add5323-fc98-4e71-a783-27dbb922c99f"
                                        }
                                    ],
                                    "isActive": true,
                                    "parentId": null,
                                    "isRequired": false,
                                    "subQuestions": []
                                }
                            }
                        ],
                        "schema": "serviceDefinition",
                        "version": 1
                    }
                },
                {
                    "id": "bd8d4fe4-ebfd-4b1e-95cd-8864d54b67cd",
                    "referenceId": "90a4fc55-c90e-404b-a45f-bfb3a03a473b",
                    "tenantId": "mz",
                    "code": "SN2.CLINICS.SN1",
                    "dataType": "SingleValueList",
                    "values": [
                        "YES",
                        "NO",
                        "NOT_SELECTED"
                    ],
                    "isActive": true,
                    "required": false,
                    "regex": null,
                    "order": "4",
                    "auditDetails": {
                        "createdBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "lastModifiedBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "createdTime": 1746526115203,
                        "lastModifiedTime": 1746526115203
                    },
                    "additionalFields": {
                        "fields": [
                            {
                                "key": "600cc1bb-3b84-48db-898c-9d80c88967b8",
                                "value": {
                                    "id": "c65ac34b-7cc0-4993-a8fe-37e854d2b189",
                                    "key": 4,
                                    "type": {
                                        "code": "SingleValueList"
                                    },
                                    "level": 2,
                                    "title": "Do you have enough products for distribution to health facilities?",
                                    "value": null,
                                    "options": [
                                        {
                                            "id": "cb45ca84-7cc0-4993-a8fe-37e854d2b189",
                                            "key": 1,
                                            "label": "Yes",
                                            "subQuestions": [],
                                            "optionDependency": false,
                                            "parentQuestionId": "c65ac34b-7cc0-4993-a8fe-37e854d2b189"
                                        },
                                        {
                                            "id": "a54c73cb-60da-4c51-8501-cf4a4f473a66",
                                            "key": 2,
                                            "label": "No",
                                            "subQuestions": [],
                                            "optionComment": true,
                                            "optionDependency": false,
                                            "parentQuestionId": "c65ac34b-7cc0-4993-a8fe-37e854d2b189"
                                        }
                                    ],
                                    "isActive": true,
                                    "parentId": "23ace43b-e0b5-428f-9d11-12fc5b10b1ac1",
                                    "isRequired": false,
                                    "subQuestions": []
                                }
                            }
                        ],
                        "schema": "serviceDefinition",
                        "version": 1
                    }
                },
                {
                    "id": "9d974565-94e1-421c-8778-d6999fdc2fc9",
                    "referenceId": "90a4fc55-c90e-404b-a45f-bfb3a03a473b",
                    "tenantId": "mz",
                    "code": "SN3",
                    "dataType": "SingleValueList",
                    "values": [
                        "MEDICAL_EQUIPMENT",
                        "PHARMACEUTICALS",
                        "PERSONAL_PROTECTIVE_EQUIPMENT_(PPE)",
                        "NOT_SELECTED"
                    ],
                    "isActive": true,
                    "required": false,
                    "regex": null,
                    "order": "3",
                    "auditDetails": {
                        "createdBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "lastModifiedBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                        "createdTime": 1746526115203,
                        "lastModifiedTime": 1746526115203
                    },
                    "additionalFields": {
                        "fields": [
                            {
                                "key": "0e60aba6-ed2f-4045-8ae7-e70b9a501f86",
                                "value": {
                                    "id": "23ca54be-038e-42df-a557-bb5fcd374dd5",
                                    "key": 3,
                                    "type": {
                                        "code": "SingleValueList"
                                    },
                                    "level": 1,
                                    "title": "What services or products do you distribute to health facilities?",
                                    "value": null,
                                    "options": [
                                        {
                                            "id": "ea32bc56-038e-42df-a557-bb5fcd374dd5",
                                            "key": 1,
                                            "label": "Medical equipment",
                                            "subQuestions": [],
                                            "optionComment": false,
                                            "optionDependency": false,
                                            "parentQuestionId": "23ca54be-038e-42df-a557-bb5fcd374dd5"
                                        },
                                        {
                                            "id": "a34vc429-d13f-4340-ae5e-fe7f8aca4212",
                                            "key": 2,
                                            "label": "Pharmaceuticals",
                                            "subQuestions": [],
                                            "optionDependency": false,
                                            "parentQuestionId": "23ca54be-038e-42df-a557-bb5fcd374dd5"
                                        },
                                        {
                                            "id": "6c43b57c-d13f-4340-ae5e-fe7f8aca4212",
                                            "key": 3,
                                            "label": "Personal protective equipment (PPE)",
                                            "subQuestions": [],
                                            "optionDependency": false,
                                            "parentQuestionId": "23ca54be-038e-42df-a557-bb5fcd374dd5"
                                        }
                                    ],
                                    "isActive": true,
                                    "parentId": null,
                                    "isRequired": false,
                                    "subQuestions": []
                                }
                            }
                        ],
                        "schema": "serviceDefinition",
                        "version": 1
                    }
                }
            ],
            "auditDetails": {
                "createdBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                "lastModifiedBy": "8ca0fd96-d0d8-4c1d-b209-4aa5518f78e7",
                "createdTime": 1746526115203,
                "lastModifiedTime": 1746526115203
            },
            "additionalFields": {
                "fields": [
                    {
                        "key": "ed91b117-88d1-46b8-9501-2ae8487cb1c5",
                        "value": {
                            "name": "WAREHOUSE DISTRIBUTOR",
                            "role": "DISTRIBUTOR",
                            "type": "WAREHOUSE",
                            "helpText": ""
                        }
                    }
                ],
                "schema": "serviceDefinition",
                "version": 1
            },
            "clientId": null
        }
    ],
    "Pagination": null
}

// Fix for CheckListConfig function
export const CheckListConfig = (values) => {
    //const { t } = useTranslation();

    const createConfig = (field, label, codes) => {
        let type = field.dataType === "SingleValueList" ? "radio" : "text";
        return {
            isMandatory: field.required,
            key: field.code,
            type: type,
            label: field.code,//t(`${label}.${codes}`),
            disable: false,
            selectedOption: values[field.code],
            populators: {
                name: field.code,
                optionsKey: "name",
                alignVertical: true,
                options: field.values?.slice(0, -1).map(item => ({
                    code: item,
                    name: item, //t(`${label}.${codes}.${item}`),
                }))
            },
        };
    };

    let config = [];
    let fields = response.ServiceDefinitions[0].attributes;
    fields.forEach(item => {
        const codeParts = item.code.split(".");
        if (codeParts.length === 1) {
            config.push(createConfig(item, response.ServiceDefinitions[0].code, item.code));
        }
        else if (codeParts.length > 1) {
            const code = codeParts[0];
            const value = codeParts[1];

            // Here's the fix: Check if values[code] is an object or a simple value
            const selectedValue = values[code]?.code || values[code];
            if (values[code] && selectedValue === value) {
                config.push(createConfig(item, response.ServiceDefinitions[0].code, item.code));
            }
        }
    });

    return [
        {
            body: config
        }
    ];
};
