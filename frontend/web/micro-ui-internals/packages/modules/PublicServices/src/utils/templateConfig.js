export const AddressFields = [
    {
      "name": "address",
      "label": "Address ",
      "type": "object",
        "properties": [
          {
            "name": "pincode",
            "label": "Pincode ",
            "disable" : false,
            "type": "string",
            "format": "number",
            "maxLength": 6,
            "minLength": 0,
            "validation": {
              "regex": "^[1-9][0-9]{5}$",
              "message": "Only 6 numbers allowed",
              "maxLength": 6,
              "minLength": 0,
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
            "format": "radioordropdown",
            "required": false,
          },
          {
            "name": "streetName",
            "label": "Street Name ",
            "disable" : false,
            "type": "string",
            "format": "text",
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
    }
  ]

export const ApplicantFields =  [{
    "name": "applicantDetails",
    "label": "Applicant Details ",
    "type": "array",
    "items":{
      "type": "object",
      "properties": [
        {
          "name": "OwnerName",
          "label": "Owner Name ",
          "disable" : false,
          "type": "string",
          "format": "text",
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
          "format": "mobileNumber",
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
          "format": "radioordropdown",
          "reference": "mdms",
          "required": false,
          "schema": "common-masters.GenderType" 
        },
      ]
    }
  }]