import React, { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
//import { searchGenericConfig  } from "../../../configs/searchGenericConfig";
import {useSearchGenericConfig} from "../../../configs/searchGenericConfig";
import { InboxSearchComposer, Loader } from "@egovernments/digit-ui-components";
import { useParams } from "react-router-dom/cjs/react-router-dom.min";

const SearchTL = () => {
    const { t } = useTranslation();
    const { module } = useParams();
    const tenantId = Digit.ULBService.getCurrentTenantId();

    const configs = useSearchGenericConfig();


    const updatedConfig = useMemo(() => Digit.Utils.preProcessMDMSConfigInboxSearch(t, configs, "sections.search.uiConfig.fields", {
        updateDependent : [
          {
            key : "businessService",
            value : [{code:"NewTL", name:"NewTL", serviceCode:"SVC-DEV-TRADELICENSE-NEWTL-04"},{code:"OldTL", name:"OldTL", serviceCode:"SVC-DEV-TRADELICENSE-OLDTL-07"}]
          }
        ]
      }), [
        configs
      ]);

    return(
    <React.Fragment>
        <div className="digit-inbox-search-wrapper">
            <InboxSearchComposer configs={updatedConfig}></InboxSearchComposer>
        </div>
    </React.Fragment>
    );
};

export default SearchTL;