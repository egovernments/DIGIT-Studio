import { FormComposerV2, Stepper, Toast } from "@egovernments/digit-ui-components";
import React, { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useHistory, useParams } from "react-router-dom";
import { serviceConfigPGR } from "../../../configs/serviceConfigurationPGR";
import { serviceConfig } from "../../../configs/serviceConfiguration";
import { generateFormConfig } from "../../../utils/generateFormConfigFromSchemaUtil";
import { transformToApplicationPayload } from "../../../utils";
import { Loader } from "@egovernments/digit-ui-react-components";

const DigitDemoComponent = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const [showToast, setShowToast] = useState(null);
  const { module } = useParams();
  const { service } = useParams();
  let serviceCode = `${module.toUpperCase()}_${service.toUpperCase()}`;

  const [currentStep, setCurrentStep] = useState(1);
  const [formData, setFormData] = useState({});
  const tenantId = Digit.ULBService.getCurrentTenantId();
  const queryStrings = Digit.Hooks.useQueryParams();

  //TODO: hardcoded the config for now, need to be changed in future

  // const requestCriteria = {
  //   url: "/egov-mdms-service/v2/_search",
  //   body: {
  //     MdmsCriteria: {
  //       tenantId: tenantId,
  //       schemaCode: "Studio.ServiceConfiguration"
  //     },
  //   },
  //   //changeQueryName: "sorOverhead"
  // };
  // const {isLoading: moduleListLoading, data} = Digit.Hooks.useCustomAPIHook(requestCriteria);

  // let config = data?.mdms?.filter((item) => item?.uniqueIdentifier.toLowerCase() === `${module}.${service}`.toLowerCase())[0];

  let config = serviceConfigPGR;

  let Updatedconfig = {
    ServiceConfiguration : [config?.ServiceConfiguration[0]],
    tenantId: tenantId,
    module: module,
  }

  // const configMap = {
  //   pgr: serviceConfigPGR,
  //   TradeLicense: serviceConfig
  // };
  // console.log(configMap[module],"configMap")

  const rawConfig = generateFormConfig(Updatedconfig, module.toUpperCase(),service?.toUpperCase());
  const steps = rawConfig.map((config) => config.head || config.label || "Untitled Section");

  const currentFormConfig = rawConfig[currentStep - 1];
  let schemaCode = queryStrings?.serviceCode || "SVC-DEV-TRADELICENSE-NEWTL-04";

  const reqCreate = {
    url: `/public-service/v1/application/${schemaCode}`,
    params: {},
    body: {},
    method: "POST",
    headers: {},
    config: {
      enable: false,
    },
  };

  const mutation = Digit.Hooks.useCustomAPIMutationHook(reqCreate);

  const onSubmit = async (data) => {
    const sectionName = currentFormConfig.name || `section_${currentStep}`;
  
    const updatedFormData = currentFormConfig?.type === "multiChildForm" || currentFormConfig?.type === "documents" ? { ...formData, ...data } : { ...formData, [sectionName]: data }; 
    setFormData(updatedFormData);

    const isLastStep = currentStep === rawConfig.length;

    if (!isLastStep) {
      setCurrentStep((prev) => prev + 1);
    } else {
      // Final submit
      await mutation.mutate(
        {
          url: `/public-service/v1/application/${schemaCode}`,
          params: {},
          headers: { "x-tenant-id": tenantId },
          method: "POST",
          body: transformToApplicationPayload(updatedFormData,Updatedconfig,service,tenantId),
          config: {
            enable: true,
          },
        },
        {
          onSuccess: (data) => {
            history.push({
              pathname: `/${window.contextPath}/employee/publicservices/${module}/${service}/response`,
              search: "?isSuccess=true",
              state: {
                message: "Application Created Successfully",
                showID: true,
                applicationNumber: data?.Application?.applicationNumber,
                redirectionUrl :  `/${window.contextPath}/employee/publicservices/${module}/${service}/ViewScreen?applicationNumber=${data?.Application?.applicationNumber}&serviceCode=${schemaCode}`,
              },
            });
          },
          onError: () => {
            history.push({
              pathname: `/${window.contextPath}/employee/publicservices/${module}/response`,
              search: "?isSuccess=false",
              state: {
                message: "Application Creation Failed",
                showID: false,
              },
            });
          },
        }
      );
    }
  };

  const onStepperClick = (stepIndex) => {
    const clickedStepIndex = stepIndex + 1; // because currentStep is 1-based
    const clickedHead = rawConfig[stepIndex].name;
    if (Object.keys(formData).includes(clickedHead)) {
      setCurrentStep(clickedStepIndex);
    }
  };

  const closeToast = () => {
    setShowToast(false);
  };


  // if (moduleListLoading) {
  //   return <Loader />;
  // }

  console.log(formData[currentFormConfig?.name || `section_${currentStep}`],"mmmmmmm")
  console.log(formData,"formdata");
  return (
    <React.Fragment>
      <Stepper
        customSteps={steps}
        currentStep={currentStep}
        onStepClick={onStepperClick}
        activeSteps={currentStep}
      />
      <FormComposerV2
        heading={t(`${serviceCode}_HEADING`)}
        label={currentStep === steps.length ? t(`${serviceCode}_SUBMIT`) : t(`${serviceCode}_NEXT`)}
        description={" "}
        text={" "}
        config={[{
          ...currentFormConfig,
          body: currentFormConfig?.body?.filter((a) => !a.hideInEmployee),
        }]}
        defaultValues={{...formData[currentFormConfig?.name || `section_${currentStep}`] || {}}}
        onSubmit={onSubmit}
        fieldStyle={{ marginRight: 0 }}
      />
      {showToast &&
        <Toast
          style={{ zIndex: "10000" }}
          error={showToast?.error}
          label={t(showToast?.message)}
          onClose={closeToast}
          isDleteBtn={true}
        />}
    </React.Fragment>
  );
};

export default DigitDemoComponent;
