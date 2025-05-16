import React from "react";
import { useLocation } from "react-router-dom/cjs/react-router-dom.min";
import { useState, useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
import { FormComposerV2, Loader } from "@egovernments/digit-ui-components";
import CheckListConfig from "../../../configs/checkListConfig.js";
import { updateCheckListConfig } from "../../../configs/checkListConfig.js";
import { useParams } from "react-router-dom";
import transformViewCheckList from "../../../utils/createUtils.js";

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
  const { accid, id, code } = useParams();
  const { t } = useTranslation();
  const [cardItems, setCardItems] = useState([]);
  const [states, dispatch] = useReducer(formReducer, {
    formData: {}
  });

  const [config, setConfig] = useState(null);

  const request = {
    url: "/health-service-request/service/definition/v1/_search",
    params: {},
    body: {},
    method: "POST",
    headers: {},
    config: {
      enable: false,
    },
  }
  const mutation = Digit.Hooks.useCustomAPIMutationHook(request);

  const getcarditems = async (code) => {
    await mutation.mutate(
      {
        url: "/health-service-request/service/definition/v1/_search",
        method: "POST",
        body: transformViewCheckList(code),
        config: {
          enable: false,
        },
      },
      {
        onSuccess: (res) => {
          console.log(res, "application_response");
          setCardItems(res?.ServiceDefinitions || []);
        },
        onError: () => {
          console.log("Error occurred");
          setCardItems([]);
        },
      }
    )
  }
  useEffect(() => {
    getcarditems([code]);
  }, [code]);

  useEffect(() => {
    if (cardItems && cardItems.length > 0) {
      setConfig(CheckListConfig(cardItems));
    }
  }, [cardItems]);

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
      {config ? (
        <FormComposerV2
          label={t("Submit")}
          config={config}
          onFormValueChange={(setValue, formData) => { handleFormValueChange(formData) }}
          onSubmit={onSubmit}
          fieldStyle={{ marginRight: 2 }}
        />
      ) : (
        <Loader/>
      )}
    </div>
  );
};

export default CheckList;