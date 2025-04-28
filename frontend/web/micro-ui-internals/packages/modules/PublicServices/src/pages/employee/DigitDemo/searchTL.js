import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
import { searchConfig } from "../../../configs/searchConfigurationTL";
import { InboxSearchComposer } from "@egovernments/digit-ui-components";

const SearchTL = () => {
    const { t } = useTranslation();

    const onSubmit = (data) => {
        console.log(data, "Final Submit Data");
    };

    const configs = searchConfig;

    return(
    <React.Fragment>
        <div className="digit-inbox-search-wrapper">
            <InboxSearchComposer configs={configs}></InboxSearchComposer>
        </div>
    </React.Fragment>
    );
};

export default SearchTL;