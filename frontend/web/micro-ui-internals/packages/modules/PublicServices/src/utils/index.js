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
  



  export const generateViewConfigFromResponse = (application, t) => {
    // const getTranslatedValue = (prefix, value) => {
    //   if (!value) return "NA";
    //   return t ? t(`${prefix}_${value}`) : value;
    // };

    const extractSectionValues = (data, prefix) => {
      const shouldTranslate = (value) => {
        if (typeof value !== "string") return false;
        const cleaned = value.replace(/-/g, "_");
        const hasOnlyNumbersOrDate = /^[\d_\-]+$/.test(value); // matches 2025-04-16, 2025_04_16, etc.
    
        return (
          cleaned.includes("_") &&
          !hasOnlyNumbersOrDate &&  // don't translate if it's numbers/dates
          /^[A-Z_]+$/.test(cleaned)  // must be UPPERCASE and underscore
        );
      };
    
      const formatField = (key, value) => {
        const isTranslate = shouldTranslate(value);
        const cleanedValue = typeof value === "string" ? value.replace(/-/g, "_") : value;
        return {
          key: t(`${application?.module.toUpperCase()}_${application?.businessService.toUpperCase()}_${key.toUpperCase()}`),
          value: isTranslate ? t(`COMMON_${key.toUpperCase()}_${cleanedValue}`) : (value || "NA"),
          isTranslate,
        };
      };
    
      if (Array.isArray(data)) {
        return data.flatMap((item, index) => {
          const itemFields = Object.keys(item || {})
            .filter((key) => {
              const value = item[key];
              return key.toLowerCase() !== "id" && value !== undefined && value !== null && value !== "";
            })
            .map((key) => formatField(key, item[key]));
    
          if (itemFields.length > 0) {
            return [
              {
                key: `Item ${index + 1}`,
                value: "",
                isTranslate: false,
              },
              ...itemFields
            ];
          }
          return [];
        });
      } else {
        return Object.keys(data || {})
          .filter((key) => {
            const value = data[key];
            return key.toLowerCase() !== "id" && value !== undefined && value !== null && value !== "";
          })
          .map((key) => formatField(key, data[key]));
      }
    };
    
    
    
  
    const serviceDetails = application.serviceDetails || {};
    const addressDetails = application.address || {};
    const applicantDetails = application.applicants?.[0] || {};
  
    const cards = [];
  
    // Service Details card
    if (Object.keys(serviceDetails).length > 0) {
      const serviceSections = Object.keys(serviceDetails)
        .map((serviceKey) => {
          const data = serviceDetails[serviceKey];
          const values = extractSectionValues(
            data,
            `${serviceKey.toUpperCase()}`
          );
          if (values.length > 0) {
            const headerKey = `${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}_${serviceKey.toUpperCase()}`;
            return {
              head: headerKey,
              type: "DATA",
              sectionHeader: { value: headerKey, inlineStyles: {} },
              values,
            };
          }
          return null;
        })
        .filter(Boolean);
  
      if (serviceSections.length > 0) {
        cards.push({
          sections: serviceSections,
        });
      }
    }
  
    // Address Details card
    const addressValues = extractSectionValues(addressDetails, "ADDRESS");
    if (addressValues.length > 0) {
      const headerKey = "ADDRESS_DETAILS";
      cards.push({
        sections: [
          {
            head: headerKey,
            type: "DATA",
            sectionHeader: { value: headerKey, inlineStyles: {} },
            values: addressValues,
          },
        ],
      });
    }
  
    // Applicant Details card
    const applicantValues = extractSectionValues(applicantDetails, "APPLICANT");
    if (applicantValues.length > 0) {
      const headerKey = "APPLICANT_DETAILS";
      cards.push({
        sections: [
          {
            head: headerKey,
            type: "DATA",
            sectionHeader: { value: headerKey, inlineStyles: {} },
            values: applicantValues,
          },
        ],
      });
    }
  
    cards.push(
      {
        navigationKey: "card1",
        sections: [
          {
            type: "WFHISTORY",
            businessService: "NewTL",
            applicationNo: "APL-DEV-TL-NEWTL-28",
            tenantId: "dev",
            timelineStatusPrefix: `WF_${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}`,
            breakLineRequired: false,
            config : {
              select: (data) => {
                return {...data, timeline: data?.timeline.filter((ob) => ob?.performedAction !== "DRAFT")}
              },
            }
          },
          {
            type: "WFACTIONS",
            forcedActionPrefix: `WF_${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}_ACTION`,
            businessService: "NewTL",
            applicationNo: "APL-DEV-TL-NEWTL-28",
            tenantId: "dev",
            applicationDetails: application,
            url: "/public-service/v1/application/SVC-DEV-TRADELICENSE-NEWTL-04",
            moduleCode: "TradeLicense",
            //editApplicationNumber: undefined,
            //editCallback : getRedirectionCallback
          },
        ],
      }
    )
    const config = {
      cards,
      apiResponse: application,
      additionalDetails: application.additionalDetails || {},
      horizontalNav: {
        showNav: false,
        configNavItems: [],
        activeByDefault: "",
      },
    };
  
    return config;
  };

  export const transformResponseforModulePage = (data) => {
    const moduleData = {}; // Object to store modules and their corresponding business services
  
    // Process each item
    data.forEach((item) => {
      const module = item.module;
  
      // If module is already processed, add the businessService to its list
      if (!moduleData[module]) {
        moduleData[module] = {
          heading: `${module.toUpperCase()}_HEADING`,
          cardDescription: `${module.toUpperCase()}_CARDDESCRIPTION`,
          businessServices: new Set(), // Set to store unique businessServices
          module: module,
          //serviceCode : item?.serviceCode
        };
      }
  
      // Add the businessService to the set (to ensure uniqueness)
      moduleData[module].businessServices.add({businessService : item.businessService, serviceCode: item?.serviceCode});
    });
  
    // Convert the moduleData object to an array of objects
    return Object.keys(moduleData).map((module) => {
      const moduleInfo = moduleData[module];
      return {
        heading: moduleInfo.heading,
        cardDescription: moduleInfo.cardDescription,
        businessServices: Array.from(moduleInfo.businessServices), // Convert the Set to an array
        module: module,
        //serviceCode : moduleInfo?.serviceCode,
      };
    });
  };
  
  
  
  
  
  


export default {};