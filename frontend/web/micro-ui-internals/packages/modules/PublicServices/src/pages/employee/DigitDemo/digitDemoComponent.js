import { FormComposerV2, Stepper, Toast } from "@egovernments/digit-ui-components";
import React, { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useHistory, useParams } from "react-router-dom";
import { serviceConfigPGR } from "../../../configs/serviceConfigurationPGR";
import { serviceConfig } from "../../../configs/serviceConfiguration";
import { generateFormConfig } from "../../../utils/generateFormConfigFromSchemaUtil";
import { transformToApplicationPayload } from "../../../utils";

const DigitDemoComponent = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const [showToast, setShowToast] = useState(null);
  const { module } = useParams();
  const { service } = useParams();

  const [currentStep, setCurrentStep] = useState(1);
  const [formData, setFormData] = useState({});
  const tenantId = Digit.ULBService.getCurrentTenantId();
  
  const configMap = {
    pgr: serviceConfigPGR,
    tl: serviceConfig
  };

  const rawConfig = generateFormConfig(configMap[module], module.toUpperCase());
  const steps = rawConfig.map((config) => config.head || config.label || "Untitled Section");

  const currentFormConfig = rawConfig[currentStep - 1];

  const reqCreate = {
    url: `/public-service/v1/application/SVC-DEV-TRADELICENSE-NEWTL-04`,
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
  
    const updatedFormData = currentFormConfig?.type === "multiChildForm" ? { ...formData, ...data } : { ...formData, [sectionName]: data };
    setFormData(updatedFormData);

    const isLastStep = currentStep === rawConfig.length;

    if (!isLastStep) {
      setCurrentStep((prev) => prev + 1);
    } else {
      // Final submit
      await mutation.mutate(
        {
          url: `/public-service/v1/application/SVC-DEV-TRADELICENSE-NEWTL-04`,
          params: {},
          headers: { "x-tenant-id": tenantId },
          method: "POST",
          body: transformToApplicationPayload(updatedFormData,configMap[module],service,tenantId),
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
    const clickedHead = rawConfig[stepIndex].head;
    if (Object.keys(formData).includes(clickedHead)) {
      setCurrentStep(clickedStepIndex);
    }
  };

  const closeToast = () => {
    setShowToast(false);
  };

  return (
    <React.Fragment>
      <Stepper
        customSteps={steps}
        currentStep={currentStep}
        onStepClick={onStepperClick}
        activeSteps={currentStep}
      />
      <FormComposerV2
        heading={t("Local Business License Issuing System")}
        label={currentStep === steps.length ? t("Submit") : t("Next")}
        description={" "}
        text={" "}
        config={[{
          ...currentFormConfig,
          body: currentFormConfig.body.filter((a) => !a.hideInEmployee),
        }]}
        defaultValues={formData[currentFormConfig.head] || {}}
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
