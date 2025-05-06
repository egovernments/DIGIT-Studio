import React, { useState, useEffect, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useHistory, useParams } from "react-router-dom";
import { searchConfig } from "../../../configs/searchConfiguration";
import { InboxSearchComposer } from "@egovernments/digit-ui-components";
import { Loader } from "@egovernments/digit-ui-react-components";

const SearchService = () => {
    const { t } = useTranslation();
    const { module } = useParams();
    const configs = searchConfig(module);
    const tenantId = Digit.ULBService.getCurrentTenantId();

    const request = {
        url: "/public-service/v1/service",
        headers: {
            "X-Tenant-Id": tenantId
        },
        method: "GET",
    }
    const { isLoading, data } = Digit.Hooks.useCustomAPIHook(request);
    const services = data?.Services.map(service => service.businessService);

    useMemo(() => {
        Digit.Utils.preProcessMDMSConfigInboxSearch(t, configs, "sections.search.uiConfig.fields", {
            updateDependent: [
                {
                    key: "service",
                    value: services?.map(item => ({
                        code: item,
                        names: item
                    }))
                },
            ]
        }
        )
    }, [services]);

    if (isLoading) {
        return <Loader />;
    }
    return (
        <React.Fragment>
            <div className="digit-inbox-search-wrapper">
                <InboxSearchComposer configs={configs}></InboxSearchComposer>
            </div>
        </React.Fragment>
    );
};

export default SearchService;