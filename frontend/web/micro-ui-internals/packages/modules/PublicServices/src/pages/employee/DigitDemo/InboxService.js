import React, { useState, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { InboxConfig } from "../../../configs/InboxConfiguration";
import { InboxSearchComposer } from "@egovernments/digit-ui-components";
import { useHistory, useParams } from "react-router-dom";

const InboxService = () => {
    const { t } = useTranslation();
    const { module } = useParams();
    const configs = InboxConfig();
    const tenantId = Digit.ULBService.getCurrentTenantId();

    // const request = {
    //     url: "/public-service/v1/service",
    //     headers: {
    //         "X-Tenant-Id": tenantId
    //     },
    //     method: "GET",
    // }
    // const { isLoading, data } = Digit.Hooks.useCustomAPIHook(request);
    // const services = data?.Services.map(service => service.businessService);

    const updateConfig = useMemo(() => {
        Digit.Utils.preProcessMDMSConfigInboxSearch(t, configs, "sections.search.uiConfig.fields", {
            updateDependent: [
                {
                    key: "businessService",
                    value: [{ code: "NewTL", name: "NewTL"}]
                },
            ]
        }
        )
    }, [configs]);

    return (
        <React.Fragment>
            <div className="digit-inbox-search-wrapper">
                <InboxSearchComposer configs={configs}></InboxSearchComposer>
            </div>
        </React.Fragment>
    );
};

export default InboxService;