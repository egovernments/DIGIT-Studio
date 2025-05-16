import _ from "lodash";
import { UICustomizations } from "../configs/UICustomizations";
import { useQuery, useQueryClient } from "react-query";

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
    const excludedKeys = ["address", "applicantDetails", "uploadedDocs", "uploaded"];
    const validSections = Object.keys(formData).reduce((acc, key) => {
      if (!excludedKeys.includes(key)) {
        acc[key] = formData[key];
      }
      return acc;
    }, {});
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
    console.log(serviceDetails,"service-details");
    return serviceDetails;
  };
  


  const transformUploadedDocs = (uploadedDocs = {}) => {
    const documents = [];
  
    Object.entries(uploadedDocs).forEach(([docType, docEntries]) => {
      docEntries?.forEach(([fileName, docMeta]) => {
        const fileStoreId = docMeta?.fileStoreId?.fileStoreId;
  
        if (fileStoreId) {
          documents.push({
            documentType: docType,
            fileStoreId: fileStoreId,
            documentUid: null, // Reusing fileStoreId as UID if no separate UID exists
            additionalDetails: {}
          });
        }
      });
    });
  
    return documents;
  };  

  export const transformToApplicationPayload = (formData, configMap, service, tenantId) => {
    const currentConfig = configMap?.ServiceConfiguration?.find(ob => ob?.service === service);
  
    const serviceDetails = getServiceDetails(formData);
  
    const applicants = formData.applicantDetails?.filter(Boolean)?.map((applicant, index) => ({
      type: "CITIZEN",
      name: applicant?.OwnerName,
      mobileNumber: Number(applicant?.mobileNumber),
      emailId: applicant?.email || `user${index + 1}@example.com`,
      prefix: "91",
      active: true,
    })) || [];
  
    const documents = transformUploadedDocs(formData?.uploadedDocs);
  
    const requestBody = {
      Application: {
        tenantId,
        module: currentConfig?.module,
        businessService: currentConfig?.service,
        status: "INACTIVE",
        channel: "counter",
        reference: null,
        workflowStatus: "applied",
        serviceDetails: {
          ...serviceDetails
        },
        applicants,
        address: {
          tenantId,
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
        documents, // <-- documents as top-level key
        additionalDetails: {
          ref1: "val1" 
        },
        Workflow: {
          action: "CREATE",
          comment: "",
          assignees: []
        }
      }
    };
  
    return requestBody;
  };



  export const generateViewConfigFromResponse = (application, t) => {
  
    const extractSectionValues = (data, prefix) => {
      const shouldTranslate = (value) => {
        if (typeof value !== "string") return false;
        const cleaned = value.replace(/-/g, "_");
        const hasOnlyNumbersOrDate = /^[\d_\-]+$/.test(value);
        return (
          cleaned.includes("_") &&
          !hasOnlyNumbersOrDate &&
          /^[A-Z_]+$/.test(cleaned)
        );
      };
  
      const formatField = (key, value) => {
        const isTranslate = shouldTranslate(value);
        const cleanedValue = typeof value === "string" ? value.replace(/-/g, "_") : value;
        return {
          key: t
            ? t(`${application?.module.toUpperCase()}_${application?.businessService.toUpperCase()}_${key.toUpperCase()}`)
            : key,
          value: isTranslate ? (t ? t(`COMMON_${key.toUpperCase()}_${cleanedValue}`) : cleanedValue) : (value || "NA"),
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
                key: `${t(`${application?.module.toUpperCase()}_${application?.businessService?.toUpperCase()}_${prefix.toUpperCase()}`)} ${index + 1}`, // <-- Label only, no value
                value: "",
                isTranslate: false,
              },
              ...itemFields,
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
    const applicants = application.applicants || [];
  
    const cards = [];
  
    // Service Details card
    if (Object.keys(serviceDetails).length > 0) {
      const serviceSections = Object.keys(serviceDetails)
        .map((serviceKey) => {
          const data = serviceDetails[serviceKey];
          const values = extractSectionValues(data, `${serviceKey.toUpperCase()}`);
          if (values.length > 0) {
            const headerKey = `${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}_${serviceKey.toUpperCase()}`;
            return {
              head: t(headerKey),
              type: "DATA",
              sectionHeader: { value: t(headerKey), inlineStyles: {} },
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
      const headerKey = `${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}_ADDRESS_DETAILS`;
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
  
    // Applicant Details card (single header, multiple applicants with labels)
    if (Array.isArray(applicants) && applicants.length > 0) {
      const applicantValues = extractSectionValues(applicants, "APPLICANT");
      cards.push({
        sections: [
          {
            head: `${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}_APPLICANT DETAILS`,
            type: "DATA",
            sectionHeader: { value: `${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}_APPLICANT DETAILS`, inlineStyles: {} },
            values: applicantValues,
          },
        ],
      });
    }
    //documents enablement
    const rawDocuments = application?.additionalDetails?.documents || {};
    const flattenedDocuments = [];

    Object.entries(rawDocuments).forEach(([docType, docEntries]) => {
      docEntries.forEach((entry) => {
        const [fileName, fileObj] = entry || [];
        const fileStoreId = fileObj?.fileStoreId?.fileStoreId;

        if (fileStoreId) {
          flattenedDocuments.push({
            title: docType || "NA",
            documentType: docType || "NA",
            documentUid: fileName || "NA",
            fileStoreId: fileStoreId,
          });
        }
      });
    });

    if (flattenedDocuments.length > 0) {
      cards.push({
        navigationKey: "card-documents",
        sections: [
          {
            type: "DOCUMENTS",
            documents: [
              {
                title: `${application?.module.toUpperCase()}_${application?.businessService.toUpperCase()}_DOCUMENTS`, // or any module-specific label
                BS: application.module || "Module",
                values: flattenedDocuments,
              },
            ],
            inlineStyles: {
              // marginTop: "1rem",
            },
          },
        ],
      });
    }
  
    // Workflow History & Actions card
    cards.push({
      navigationKey: "card1",
      sections: [
        {
          type: "WFHISTORY",
          businessService: application.businessService,
          applicationNo: application.applicationNumber,
          tenantId: application.tenantId,
          timelineStatusPrefix: `WF_${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}`,
          breakLineRequired: false,
          config: {
            select: (data) => {
              return { ...data, timeline: data?.timeline?.filter((ob) => ob?.performedAction !== "DRAFT") };
            },
          },
        },
        // {
        //   type: "WFACTIONS",
        //   forcedActionPrefix: `WF_${application?.module?.toUpperCase()}_${application?.businessService?.toUpperCase()}_ACTION`,
        //   businessService: application.businessService,
        //   applicationNo: application.applicationNumber,
        //   tenantId: application.tenantId,
        //   applicationDetails: application,
        //   url: `/public-service/v1/application/${application?.serviceCode}`,
        //   moduleCode: application.module,
        // },
      ],
    });
  
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
  

  export const getServicesOptions = (services,module) => {
    const options = services?.filter((ob) => ob?.module === module && ob?.status === "ACTIVE").map((ob) =>  {return { code: ob?.businessService, name: ob?.businessService, serviceCode: ob?.serviceCode }});
    return options;
  }

  export const useWorkflowDetailsWorks = ({ tenantId, id, moduleCode, role = "CITIZEN", serviceData = {}, getStaleData, getTripData = false, config }) => {
    const queryClient = useQueryClient();
  
    const staleDataConfig = { staleTime: Infinity };
    
  
    const { isLoading, error, isError, data } = useQuery(
        ["workFlowDetailsWorks", tenantId, id, moduleCode, role, config],
        () => Digit.WorkflowService.getDetailsByIdWorks({ tenantId, id, moduleCode, role, getTripData }),
        getStaleData ? { ...staleDataConfig, ...config } : config
    );
  
    if (getStaleData) return { isLoading, error, isError, data };
  
    return { isLoading, error, isError, data, revalidate: () => queryClient.invalidateQueries(["workFlowDetailsWorks", tenantId, id, moduleCode, role]) };
  };
  
  
  
  
  


export default {};