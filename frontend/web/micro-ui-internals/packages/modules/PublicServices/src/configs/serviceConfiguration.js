export const serviceConfig = {
    "tenantId": "pg",
    "moduleName": "CommonService",
    "ServiceConfiguration": [
      {
        "service": "TL",
        "fields": [
          {
            "name": "financialYear",
            "label": "Financial Year ",
            "disable" : false,
            "type": "string",
            "defaultValue" : "2025_26",
            "prefix": "FINANCIALYEAR",
            //"reference": "mdms",
            "required": false,
            //"schema": "egf-master.FinancialYear" 
          },
          {
            "name": "licenseType",
            "label": "License Type ",
            "disable" : false,
            "defaultValue" : "PERMANENT",
            "prefix": "LICENSETYPE",
            "type": "string",
            "required": false,
          },
          {
            "name": "tradeName",
            "label": "Trade Name ",
            "disable" : false,
            "type": "string",
            "maxLength": 128,
            "minLength": 2,
            "validation": {
              "regex": "^[A-Za-z0-9 ]+$",
              "message": "Only letters and numbers allowed"
            },
            "required": false,
            "orderNumber": 1
          },
          {
            "name": "tradeStructureType",
            "label": "Trade Structure Type ",
            "disable" : false,
            "type": "string",
            "reference": "mdms",
            "required": false,
            "schema": "Tradelicence.StructureType" 
          },
          {
            "name": "tradeStructureSubType",
            "label": "Trade Structure Sub Type ",
            //same master how to give custom output for options
            // 1st way is to separate out master
            // 2nd way to pass custom logck in mdms
            "disable" : false,
            "type": "string",
            "reference": "mdms",
            "dependencies": [
              "tradeStructureType"
            ],
            "required": false,
            "schema": "Tradelicence.StructureSubType"
          },
          {
            "name": "tradeCommencementDate",
            "label": "Trade Commencement Date ",
            "disable" : false,
            "type": "date",
            "required": false,
          },
          // {
          //   "name": "tradeLicenceType",
          //   "label": "Trade Licence Type ",
          //   "disable" : false,
          //   "type": "string",
          //   "reference": "mdms",
          //   "required": true,
          //   "schema": "Tradelicence.TradeLicenceType"
          // },
          // {
          //   "name": "tradeCategory",
          //   "label": "Trade Category ",
          //   "disable" : false,
          //   "type": "string",
          //   "reference": "mdms",
          //   "required": true,
          //   "schema": "Tradelicence.TradeCategory"
          // },
          // {
          //   "name": "tradeType",
          //   "label": "Trade Type ",
          //   "disable" : false,
          //   "type": "string",
          //   "reference": "mdms",
          //   "required": true,
          //   "schema": "Tradelicence.TradeType"
          // },
          // {
          //   "name": "tradeSubType",
          //   "label": "Trade Sub Type ",
          //   "disable" : false,
          //   "type": "string",
          //   "reference": "mdms",
          //   "dependencies": [
          //     "tradeType"
          //   ],
          //   "required": true,
          //   "schema": "Tradelicence.TradeSubType"
          // },
          {
            "name": "tradeUnits",
            "label": "Trade Units ",
            "type": "object",
              "properties": [
                {
                  "name": "tradeCategory",
                  "label": "Trade Category ",
                  "disable" : false,
                  "type": "string",
                  "reference": "mdms",
                  "required": false,
                  "schema": "Tradelicence.TradeCategory"
                },
                {
                  "name": "tradeType",
                  "label": "Trade Type ",
                  "disable" : false,
                  "type": "string",
                  "reference": "mdms",
                  "required": false,
                  "schema": "Tradelicence.TradeType"
                },
                {
                  "name": "tradeSubType",
                  "label": "Trade Sub Type ",
                  "disable" : false,
                  "type": "string",
                  "reference": "mdms",
                  "required": false,
                  "schema": "Tradelicence.TradeSubType"
                }
              ]
          },
          {
            "name": "tradeAddress",
            "label": "Trade Address ",
            "type": "object",
              "properties": [
                {
                  "name": "pincode",
                  "label": "Pincode ",
                  "disable" : false,
                  "type": "string",
                  "maxLength": 6,
                  "minLength": 0,
                  "validation": {
                    "regex": "^[1-9][0-9]{5}$",
                    "message": "Only 6 numbers allowed"
                  },
                  "required": false,
                  "orderNumber": 1
                },
                {
                  "name": "city",
                  "label": "City ",
                  "disable" : false,
                  "defaultValue" : "DEV",
                  "prefix": "CITY",
                  "type": "string",
                  "required": false,
                },
                {
                  "name": "streetName",
                  "label": "Street Name ",
                  "disable" : false,
                  "type": "string",
                  "maxLength": 256,
                  "minLength": 0,
                  "validation": {
                    "regex": "^[1-9][0-9]{5}$",
                    "message": "Only 6 numbers allowed"
                  },
                  "required": false,
                  "orderNumber": 1
                },
              ]
          },
          {
            "name": "ownershipDetails",
            "label": "Ownership Details ",
            "type": "object",
              "properties": [
                {
                  "name": "OwnerName",
                  "label": "Owner Name ",
                  "disable" : false,
                  "type": "string",
                  "maxLength": 256,
                  "minLength": 0,
                  "validation": {
                    "regex": "^[1-9][0-9]{5}$",
                    "message": "Only 6 numbers allowed"
                  },
                  "required": false,
                  "orderNumber": 1
                },
                {
                  "name": "mobileNumber",
                  "label": "Mobile Number ",
                  "disable" : false,
                  "type": "mobileNumber",
                  "maxLength": 256,
                  "minLength": 0,
                  "validation": {
                    "regex": "^[6-9]\d{9}$",
                    "message": "Only 9 numbers allowed"
                  },
                  "required": false,
                  "orderNumber": 1
                },
                {
                  "name": "gender",
                  "label": "Gender ",
                  "disable" : false,
                  "type": "string",
                  "reference": "mdms",
                  "required": false,
                  "schema": "common-masters.GenderType" 
                },
              ]
          },
          {
            "name": "accessories",
            "label": "Trade accessories ",
            "type": "array",
            "items": {
              "type": "object",
              "properties": [
                {
                  "name": "accessoryType",
                  "label": "Accessory type ",
                  "disable" : false,
                  "type": "string",
                  "reference": "mdms",
                  "required": true,
                  "schema": "TradeLicense.TradeType"
                }
                // {
                //   "name": "count",
                //   "label": "accessories count ",
                //   "type": "number"
                // }
              ]
            }
          }
        ],
        "workflow": {
          "businessService": "NewTL",
          "ACTIVE": [
            "APPROVED"
          ],
          "INACTIVE": [
            "REJECTED",
            "WITHDRAWN"
          ]
        },
        "calculator": {
          "billingSlabs": [
            {
              "key": "applicationFee",
              "value": 2000
            }
          ]
        },
        "idgen": {
          "format": "tl.application.number"
        },
        "localization": {
          "modules": [
            "digit-tradelicence"
          ]
        },
        "notification": {
          "sms": {
            "TODO": "will fill later"
          },
          "email": {
            "TODO": "will fill later"
          }
        },
        "access": {
          "roles": [
            "TL_CREATOR"
          ],
          "actions": [
            {
              "url": "tl-services/v1/create"
            }
          ]
        },
        "rules": {
          "validation": {
            "type": "schema||api||custom||none",
            "service": "tradelicence",
            "schemaCode": "tradelicence.apply",
            "customFunction": "eitherhookname||function"
          },
          "calculator": {
            "type": "api||custom||none",
            "service": "tradelicence",
            "customFunction": "eitherhookname||function"
          },
          "registry": {
            "type": "api||none",
            "service": "tradelicence"
          },
          "references": [
            {
              "type": "initiate",
              "service": "tradelicence"
            }
          ]
        },
        "documents": [
          {
            "category": "address-proof",
            "documentTypes": [
              "aadhar",
              "voter"
            ],
            "active": true,
            "isMandatory": false,
            "allowedFileTypes": [
              "pdf",
              "doc",
              "docx",
              "xlsx",
              "xls",
              "jpeg",
              "jpg",
              "png"
            ],
            "maxSizeInMB": 5,
            "maxFilesAllowed": 1
          },
          {
            "category": "identity-proof",
            "documentTypes": [
              "aadhar",
              "voter"
            ],
            "active": true,
            "isMandatory": true,
            "allowedFileTypes": [
              "pdf",
              "doc",
              "docx",
              "xlsx",
              "xls",
              "jpeg",
              "jpg",
              "png"
            ],
            "maxSizeInMB": 5,
            "maxFilesAllowed": 1
          },
          {
            "category": "owner-photo",
            "documentTypes": [
              "photo"
            ],
            "active": true,
            "isMandatory": true,
            "allowedFileTypes": [
              "jpeg",
              "jpg",
              "png"
            ],
            "maxSizeInMB": 5,
            "maxFilesAllowed": 1
          }
        ],
        "pdf": [
          {
            "key": "tl-application",
            "type": "application"
          },
          {
            "key": "tl-bill",
            "type": "bill"
          },
          {
            "key": "tl-receipt",
            "type": "receipt"
          }
        ],
        "bill": {
          "service": "ApplicationFee"
        },
        "payment": {
          "gateway": "TODO"
        },
        "apiconfig": [
          {
            "type": "register",
            "host": "https://staging.digit.org||http://tl-services.egov:8080",
            "endpoint": "/tl-services/v1/create",
            "method": "post",
            "service": "tradelicence"
          },
          {
            "type": "register||calculate||validate||authenticate",
            "host": "https://staging.digit.org||http://tl-services.egov:8080",
            "endpoint": "/tl-services/v1/search",
            "method": "post",
            "service": "tradelicence"
          }
        ],
        "applicant": {
          "minimum": 1,
          "maximum": 3,
          "types": [
            "individual",
            "organisation"
          ]
        },
        "boundary": {
          "hierarchyType": "REVENUE",
          "lowestLevel": "locality"
        },
        "enabled": [
          "citizen",
          "employee"
        ]
      },
      {
        "service": "Tradelicence",
        "fields": null,
        "workflow": {
          "businessService": "ApprovedTL"
        },
        "idgen": {
          "format": "tl.licence.number"
        },
        "rules": [
          {
            "referenceType": "NewTL"
          }
        ],
        "documents": null,
        "pdf": [
          {
            "key": "tl-certificate",
            "type": "certificate"
          }
        ],
        "bill": null
      },
      {
        "service": "Tradelicence",
        "fields": null,
        "workflow": {
          "businessService": "RenewTL"
        },
        "idgen": {
          "format": "tl.renewalapplication.number"
        },
        "rules": [
          {
            "referenceType": "ApprovedTL"
          }
        ],
        "documents": [
          {
            "category": "owner-photo",
            "documentTypes": [
              "photo"
            ],
            "active": true,
            "isMandatory": false,
            "allowedFileTypes": [
              "jpeg",
              "jpg",
              "png"
            ],
            "maxSizeInMB": 5,
            "maxFilesAllowed": 1
          }
        ],
        "pdf": [
          {
            "key": "tl-renew-application",
            "type": "application"
          }
        ],
        "calculator": {
          "billingSlabs": [
            {
              "key": "applicationFee",
              "value": 2000
            }
          ]
        },
        "bill": {
          "service": "RenewalFee"
        }
      }
    ]
  }