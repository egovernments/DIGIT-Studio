import { Header, ActionBar, SubmitBar, Menu, Card, Loader, ViewComposer, MultiLink } from "@egovernments/digit-ui-react-components";
import React, { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
import { generateViewConfigFromResponse } from "../../../utils";

const DigitDemoViewComponent = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const tenantId = Digit.ULBService.getCurrentTenantId();

  const userInfo = Digit.UserService.getUser();
  const userRoles = userInfo?.info?.roles?.map((roleData) => roleData?.code);
  const [current, setCurrent] = useState(Date.now());
  const queryStrings = Digit.Hooks.useQueryParams();

  const request = {
    url : "/public-service/v1/application/SVC-DEV-TRADELICENSE-NEWTL-04",
    headers: {
      "X-Tenant-Id" : tenantId
    },
    method: "GET",
    params: {
      "ids": queryStrings?.id,
    }
  }
  const {isLoading, data} = Digit.Hooks.useCustomAPIHook(request);

  let response =  data ? data?.Application?.[0] : {};
  let config = generateViewConfigFromResponse(response,t);

  if (isLoading) {
    return <Loader />;
  }

  return (
    <React.Fragment>
      {
        <div className={"employee-application-details"} style={{ marginBottom: "24px" }}>
          {
            <Header className="works-header-view" styles={{ marginLeft: "0px", paddingTop: "10px" }}>
              {t(`${response?.module.toUpperCase()}_${response?.businessService?.toUpperCase()}_APPLICATION_DETAILS`)}
            </Header>
          }
        </div>
      }
      <ViewComposer data={config} isLoading={false} />
      {/* <ActionBar>
        {displayMenu ? <Menu localeKeyPrefix={"WORKS"} options={actionULB} optionKey={"name"} t={t} onSelect={onActionSelect} /> : null}
        <SubmitBar label={t("ACTIONS")} onSubmit={() => setDisplayMenu(!displayMenu)} />
      </ActionBar> */}
    </React.Fragment>
  );
};

export default DigitDemoViewComponent;
