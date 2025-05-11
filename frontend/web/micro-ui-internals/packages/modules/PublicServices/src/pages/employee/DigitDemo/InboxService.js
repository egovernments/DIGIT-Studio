import React, { useState, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { InboxConfig } from "../../../configs/inboxGenericConfig";
import { InboxSearchComposer } from "@egovernments/digit-ui-components";
import { useHistory, useParams } from "react-router-dom";

const InboxService = () => {
    const { t } = useTranslation();
    const { module } = useParams();
    const tenantId = Digit.ULBService.getCurrentTenantId();

    const configs=InboxConfig();

    const updatedConfig = useMemo(() => Digit.Utils.preProcessMDMSConfigInboxSearch(t, configs, "sections.search.uiConfig.fields", {
            updateDependent : [
              {
                key : "businessService",
                value : [{code:"NewTL", name:"NewTL"}]
              }
            ]
          }), [
            configs
          ]);

    return (
        <React.Fragment>
            <div className="digit-inbox-search-wrapper">
                <InboxSearchComposer configs={updatedConfig}></InboxSearchComposer>
            </div>
        </React.Fragment>
    );
};

export default InboxService;