import React from "react";
import { useLocation } from "react-router-dom/cjs/react-router-dom.min";
import { useState, useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
import { FormComposerV2 } from "@egovernments/digit-ui-components";
import CheckListConfig from "../../../configs/checkListConfig.js";
import { updateCheckListConfig } from "../../../configs/checkListConfig.js";

const formReducer = (states, action) => {
  switch (action.type) {
    case 'UPDATE_FORM':
      return {
        ...states,
        formData: action.payload
      };
    default:
      return states;
  }
};

const CheckList = () => {
  const { state } = useLocation();
  const { accountid, config: routeConfig } = state || {};
  const { t } = useTranslation();
  const [states, dispatch] = useReducer(formReducer, {
    formData: {}
  });

  const [config, setConfig] = useState(null);

  useEffect(() => {
    setConfig(CheckListConfig(routeConfig));
  }, [routeConfig]);

  const onSubmit = async (data) => {
    console.log(data, "data");
  };

  const handleFormValueChange = (formData) => {
    console.log(formData,"formdata");
    if (JSON.stringify(formData) !== JSON.stringify(states.formData)) {
      dispatch({
        type: 'UPDATE_FORM',
        payload: formData
      });
      setConfig(updateCheckListConfig(config, formData));
    }
  };

  return (
    <div>
      <FormComposerV2
        label={t("Submit")}
        config={config}
        onFormValueChange={(setValue, formData) => { handleFormValueChange(formData) }}
        onSubmit={onSubmit}
        fieldStyle={{ marginRight: 2 }}
      />
    </div>
  );
};

export default CheckList;