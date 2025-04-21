
export const generateFormConfig = (serviceConfig, serviceName) => {

  let config = [];
  if(serviceName)
  config = serviceConfig.ServiceConfiguration.find(
    (svc) => svc.service === serviceName && Array.isArray(svc.fields)
  );
  else
  config = serviceConfig;

  console.log(config,"config");
  if (!config) return [];

  let body = [];
  let fields = serviceName ? config.fields : (config?.items ? config.items.properties : config.properties);
  body = fields.map((field) => {
    let fieldType = "component"; // default
    let populators = {
      name: field.name,
    };

    if (field.type === "string" && (field.reference === "mdms" || field.defaultValue)) {
      fieldType = "radioordropdown";
      populators = {
        ...populators,
        optionsKey: "name",
        error: "sample required message",
        required: !!field.required,
        ...(
          field?.schema
            ? {
                mdmsConfig: {
                  masterName: field.schema.split(".")[1] || "Master",
                  moduleName: field.schema.split(".")[0] || "common-masters",
                  localePrefix: `COMMON_${field.name.toUpperCase()}`,
                }
              }
            : {}
        ),
        ...(
          field?.defaultValue
            ? {
              options: [
                {
                  code: field.defaultValue,
                  name: `TRADELICENSE_LICENSETYPE_${field.defaultValue}`,
                },
              ]
              }
            : {}
        ),

        // mdmsv2: {
        //             schemaCode: field?.schema,
        //           },
      };
    } else if (field.type === "string") {
      fieldType = "text";
      populators = {
        ...populators,
        error: "sample error message",
      };

      if (field.validation?.regex) {
        populators.validation = {
          pattern: new RegExp(field.validation.regex),
        };
      }
    }
    else if (field.type === "mobileNumber") {
      fieldType = "mobileNumber";
      populators = {
        ...populators,
        error: "sample error message",
      };

      if (field.validation?.regex) {
        populators.validation = {
          pattern: new RegExp(field.validation.regex),
        };
      }
    }
    else if (field.type === "date") {
      fieldType = "date";
      populators = {
        ...populators,
        error: "sample error message",
      };

      if (field.validation?.regex) {
        populators.validation = {
          pattern: new RegExp(field.validation.regex),
        };
      }
    }
    else if (field.type === "object") {
      fieldType = "childForm";
      populators = {
        ...populators,
        childform: generateFormConfig(field),
        error: "sample error message",
      };

      if (field.validation?.regex) {
        populators.validation = {
          pattern: new RegExp(field.validation.regex),
        };
      }
    }
    else if (field.type === "array") {
      fieldType = "multiChildForm";
      populators = {
        ...populators,
        childform: generateFormConfig(field),
        error: "sample error message",
      };

      if (field.validation?.regex) {
        populators.validation = {
          pattern: new RegExp(field.validation.regex),
        };
      }
    }

    return {
      label: field.label?.trim() || field.name,
      isMandatory: !!field.required,
      //description: `${field.label?.trim() || field.name} if any`,
      key: field.name,
      type: fieldType,
      disable: field?.disable || false,
      populators,
    };
  });
 

  const formconfig = [{
    head: "TradeDetails",
    subHead: " ",
    body: body,
  }]
  console.log(formconfig,"formmm")
  return formconfig;
}
