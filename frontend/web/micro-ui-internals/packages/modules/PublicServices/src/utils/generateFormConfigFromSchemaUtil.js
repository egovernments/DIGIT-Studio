import { AddressFields } from "./templateConfig";
import { ApplicantFields } from "./templateConfig";
export const generateFormConfig = (config,module,service) => {
  const serviceFields = config?.ServiceConfiguration?.[0]?.fields || [];

  const createField = (field) => {
    return {
      type: field.format || field.type,
      label : `${module}_${service}_${field.name.toUpperCase()}`,
      populators: {
        name: field.name,
        optionsKey: "name",
        error: field?.validation?.message || "field is required",
        required: !!field.required,
        validation: field.validation,
        required: field.required,
        disable: field.disable,
        defaultValue: field.defaultValue,
        prefix: field.prefix,
        reference: field.reference,
        dependencies: field.dependencies,
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
                  name: `TRADELICENSE_${field.prefix}_${field.defaultValue}`,
                },
              ]
              }
            : {}
        ),
      },
    };
  };

  const createChildForm = (objectField) => {
    return {
      head: `${module}_${service}_${objectField.name.toUpperCase()}`,
      name: objectField.name,
      body: objectField.properties.map((subField) => createField(subField)),
      type: "childform",
      step: 1,
    };
  };

  const createMultiChildForm = (arrayField) => {
    return {
      head: `${module}_${service}_${arrayField.name.toUpperCase()}`,
      name: arrayField.name,
      type: "multiChildForm",
      body: arrayField.items.properties.map((subField) => createField(subField)),
      step:2
    };
  };

  const basicFields = [];
  const stepForms = [];

  serviceFields.forEach((field) => {
    if (field.type === "object") {
      stepForms.push(createChildForm(field));
    } else if (field.type === "array") {
      stepForms.push(createMultiChildForm(field));
    } else {
      basicFields.push(createField(field));
    }
  });

  
  const addressFieldsStep = AddressFields[0]?.type === "object" ? createChildForm(AddressFields[0]) : createMultiChildForm(AddressFields[0]);
  const applicantFieldsStep = ApplicantFields[0].type === "array"? createMultiChildForm(ApplicantFields[0]) : createChildForm(ApplicantFields[0]);

  const steps = [];

  if (basicFields.length > 0) {
    steps.push({
      head: "Service Details",
      body: basicFields,
      type: "form",
    });
  }
  //need to add condition here for address
   

  return [...steps, ...stepForms, addressFieldsStep, applicantFieldsStep];

};

