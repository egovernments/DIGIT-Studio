import { FormComposerV2  } from "@egovernments/digit-ui-components";
import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
//import { serviceConfig } from "./serviceConfiguration";
import { useParams } from "react-router-dom/cjs/react-router-dom.min";
import { serviceConfigPGR } from "../../../configs/serviceConfigurationPGR";
import { serviceConfig } from "../../../configs/serviceConfiguration";
import { generateFormConfig } from "../../../utils/generateFormConfigFromSchemaUtil";

const DigitDemoComponent = () => {
  const { t } = useTranslation();

  const onSubmit = (data) => {
    ///
    console.log(data, "Final Submit Data");
  };
  const { module } = useParams();

  /* use newConfig instead of commonFields for local development in case needed */
let configMap = {
  pgr : serviceConfigPGR,
  tl : serviceConfig
}

  const configs = generateFormConfig(configMap[module],module.toUpperCase());
  return (
    <FormComposerV2
      heading={t("New Trade License Application")}
      label={t("Submit Bar")}
      description={" "}
      text={" "}
      config={configs.map((config) => {
        return {
          ...config,
          body: config.body.filter((a) => !a.hideInEmployee),
        };
      })}
      defaultValues={{}}
      onSubmit={onSubmit}
      fieldStyle={{ marginRight: 0 }}
    />
  );
};

export default DigitDemoComponent;