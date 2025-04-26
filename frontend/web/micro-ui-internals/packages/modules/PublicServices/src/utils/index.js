import _ from "lodash";
import { UICustomizations } from "../configs/UICustomizations";

  /* To Overide any existing libraries  we need to use similar method */
  const setupLibraries = (Library, service, method) => {
    window.Digit = window.Digit || {};
    window.Digit[Library] = window.Digit[Library] || {};
    window.Digit[Library][service] = method;
  };
  
  /* To Overide any existing config/middlewares  we need to use similar method */
 export const updateCustomConfigs = () => {
    setupLibraries("Customizations", "commonUiConfig", { ...window?.Digit?.Customizations?.commonUiConfig, ...UICustomizations });
    // setupLibraries("Utils", "parsingUtils", { ...window?.Digit?.Utils?.parsingUtils, ...parsingUtils });
  };

  const getServiceDetails = (formData) => {
    const { address, applicantDetails, ...validSections } = formData;
  
    const flattenValues = (obj) => {
      const flat = {};
      for (const [key, val] of Object.entries(obj)) {
        if (val && typeof val === "object" && !Array.isArray(val)) {
          flat[key] = val && typeof val === "object" && "code" in val ? val.code : val;
        } else {
          flat[key] = val;
        }
      }
      return flat;
    };
  
    const serviceDetails = {};
  
    for (const [sectionKey, sectionVal] of Object.entries(validSections)) {
      if (Array.isArray(sectionVal)) {
        // Direct arrays (not common in your example, but for safety)
        serviceDetails[sectionKey] = sectionVal.map((item) => flattenValues(item));
      } else if (typeof sectionVal === "object" && sectionVal !== null) {
        const innerKeys = Object.keys(sectionVal);
        if (innerKeys.length === 1 && Array.isArray(sectionVal[innerKeys[0]])) {
          // e.g., accessories: { accessories: [ { accessoryType: {...} } ] }
          const innerKey = innerKeys[0];
          serviceDetails[sectionKey] = {
            [innerKey]: sectionVal[innerKey].map((item) => {
              const itemKey = Object.keys(item)[0];
              const itemVal = item[itemKey];
              return {
                [itemKey]: typeof itemVal === "object" && itemVal?.code ? itemVal.code : itemVal
              };
            })
          };
        } else {
          // Normal object: flatten one level
          serviceDetails[sectionKey] = flattenValues(sectionVal);
        }
      } else {
        // Primitive value directly (unexpected case)
        serviceDetails[sectionKey] = sectionVal;
      }
    }
  
    return serviceDetails;
  };
  

  export const transformToApplicationPayload = (formData,configMap, service, tenantId) => {
   let currentConfig = configMap?.ServiceConfiguration?.filter((ob) => ob?.service === service)[0];

   let serviceDetails = getServiceDetails(formData);

    let requestBody = {
      Application: {
        tenantId: "dev",
        module: currentConfig?.module,
        businessService: currentConfig?.service,
        status: "INACTIVE",
        channel: "counter",
        reference: null,
        workflowStatus: "applied",
        serviceDetails: {
          ...serviceDetails
          // tradeName: formData.tradeName,
          // licenseType: formData.licenseType?.code,
          // tradeStructureType: formData.tradeStructureType?.code,
          // tradeStructureSubType: formData.tradeStructureSubType?.code,
          // tradeCommencementDate: formData.tradeCommencementDate,
          // tradeUnits: formData?.tradeUnits,
          //accessories: formData?.accessories?.filter(Boolean), // remove nulls
          //financialYear: formData.financialYear
        },
        applicants: [
          {
            type: "CITIZEN",
            name: formData.applicantDetails?.[0]?.OwnerName,
            userId: "2", // You can replace this dynamically
            mobileNumber: Number(formData.applicantDetails?.[0]?.mobileNumber),
            emailId: "john@gmail.com", // Optional or dynamic
            prefix: "91", // or dynamically detect
            active: true,
          }
        ],
        address: {
          tenantId: tenantId,
          latitude: 0,
          longitude: 0,
          addressNumber: "1",
          addressLine1: formData.tradeAddress?.streetName || "",
          addressLine2: "",
          landmark: "",
          city: formData.tradeAddress?.city?.name || "",
          pincode: formData.tradeAddress?.pincode,
          hierarchyType: currentConfig?.boundary?.hierarchyType,
          boundarylevel: currentConfig?.boundary?.lowestLevel,
          boundarycode: `dev.${formData.tradeAddress?.city?.code?.toLowerCase() || "city"}`,
        },
        additionalDetails: {
          ref1: "val1" // static or populate based on need
        },
        Workflow: {
          action: "APPLY",
          comment: "Applied new Application",
          assignees: []
        }
      }
    };
    return requestBody;
  }
  



export default {};