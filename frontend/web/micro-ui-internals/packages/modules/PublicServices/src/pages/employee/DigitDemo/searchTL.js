import React, { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
import { searchGenericConfig  } from "../../../configs/searchGenericConfig";
import { InboxSearchComposer, Loader } from "@egovernments/digit-ui-components";
import { getServicesOptions } from "../../../utils";
import { useParams } from "react-router-dom/cjs/react-router-dom.min";

const SearchTL = () => {
    const { t } = useTranslation();
    const { module } = useParams();
    const tenantId = Digit.ULBService.getCurrentTenantId();
    const onSubmit = (data) => {
        console.log(data, "Final Submit Data");
    };

    const configs = searchGenericConfig;

    // const request = {
    //     url : "/public-service/v1/service",
    //     headers: {
    //       "X-Tenant-Id" : tenantId
    //     },
    //     method: "GET",
    //   }
    //   const {isLoading, data} = Digit.Hooks.useCustomAPIHook(request);
    //   console.log(data,"kkoma")

    //   const memoizedData = useMemo(() => data, [JSON.stringify(data)]);

    //   const businessServiceOptions = useMemo(() => {
    //     return getServicesOptions(data?.Services, module);
    //   }, [memoizedData, module]);


    const updatedConfig = useMemo(() => Digit.Utils.preProcessMDMSConfigInboxSearch(t, configs, "sections.search.uiConfig.fields", {
        updateDependent : [
          {
            key : "businessService",
            value : [{code:"NewTL", name:"NewTL", serviceCode:"SVC-DEV-TRADELICENSE-NEWTL-04"}]
          }
        ]
      }), [
        configs
      ]);


    //   if (isLoading) {
    //     return <Loader />;
    //   }

    return(
    <React.Fragment>
        <div className="digit-inbox-search-wrapper">
            <InboxSearchComposer configs={updatedConfig}></InboxSearchComposer>
        </div>
    </React.Fragment>
    );
};

export default SearchTL;