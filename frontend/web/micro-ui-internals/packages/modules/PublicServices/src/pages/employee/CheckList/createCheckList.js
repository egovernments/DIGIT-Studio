import React from "react";
import { useState, useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
import { FormComposerV2, Loader, Toast } from "@egovernments/digit-ui-components";
import CreateCheckListConfig from "../../../configs/createCheckListConfig.js";
import { updateCheckListConfig } from "../../../configs/createCheckListConfig.js";
import { useParams } from "react-router-dom";
import transformViewCheckList from "../../../utils/createUtils.js";
import { transformCreateCheckList } from "../../../utils/createUtils.js";

const CreateCheckList = () => {
  const { accid, id, code } = useParams();
  const { t } = useTranslation();
  const [cardItems, setCardItems] = useState([]);
  const [formData, setFormData] = useState({});
  const [showToast,setShowToast] = useState(null)

  const [config, setConfig] = useState(null);

  const closeToast = () => {
    setTimeout(() => {
      setShowToast(null)
    }, 5000);
  }
 
  setTimeout(() => {
    setShowToast(null);
  }, 20000);
    

  const search_request = {
    url: "/health-service-request/service/definition/v1/_search",
    params: {},
    body: {},
    method: "POST",
    headers: {},
    config: {
      enable: false,
    },
  }
  const smutation = Digit.Hooks.useCustomAPIMutationHook(search_request);

  //application creation request
  const create_request = {
    url: "/health-service-request/service/v1/_create",
    params: {},
    body: {},
    method: "POST",
    headers: {},
    config: {
      enable: false,
    },
  }
  const cmutation = Digit.Hooks.useCustomAPIMutationHook(create_request);

  const getcarditems = async (code) => {
    await smutation.mutate(
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
      setConfig(CreateCheckListConfig(cardItems));
    }
  }, [cardItems]);

  const onSubmit = async (data) => {
    console.log(data, "data");
    const fetchdata = async (data) => {
      await cmutation.mutate(
        {
          url: "/health-service-request/service/v1/_create",
          method: "POST",
          body: transformCreateCheckList(id, accid, data),
          config: {
            enable: false,
          },
        },
        {
          onSuccess: (res) => {
            console.log(res, "application_response");
            setShowToast({ label: Digit.Utils.locale.getTransformedLocale(`${code?.replaceAll(".","_").toUpperCase()}_CREATE_SUCCESS_CHECKLIST`) })
            setCardItems(res?.ServiceDefinitions || []);
            setTimeout(() => {
              window.history.back();
            }, 3000);
          },
          onError: () => {
            console.log("Error occurred");
            setCardItems([]);
          },
        }
      )
    }
    fetchdata(data);
  };

  const handleFormValueChange = (updatedFormData) => {
    if (JSON.stringify(updatedFormData) !== JSON.stringify(formData)) {
      setFormData(updatedFormData);
      setConfig(updateCheckListConfig(config, updatedFormData));
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
        <Loader />
      )}
       {showToast && (
        <Toast
          type={showToast?.type}
          label={t(showToast?.label)}
          onClose={() => {
            setShowToast(null);
          }}
          isDleteBtn={showToast?.isDleteBtn}
        />
      )}
    </div>
  );
};

export default CreateCheckList;