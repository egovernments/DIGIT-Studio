import { AppContainer, PrivateRoute } from "@egovernments/digit-ui-react-components";
import { BreadCrumb } from "@egovernments/digit-ui-components";
import React from "react";
import { useTranslation } from "react-i18next";
import { Switch } from "react-router-dom";
import DigitDemoComponent from "./DigitDemo/digitDemoComponent";
import SearchTL from "./DigitDemo/searchTL";
import Response from "./Response";
import DigitDemoViewComponent from "./DigitDemo/digitDemoViewComponent";
import ModulePageComponent from "./DigitDemo/modulePageComponent";


const SampleBreadCrumbs = ({ location }) => {
  const { t } = useTranslation();

  const crumbs = [
    {
      internalLink: `/${window?.contextPath}/employee`,
      content: t("HOME"),
      show: true,
    },
    {
      content: t(location.pathname.split("/").pop()),
      show: true,
    },
    {
      content: t(location.pathname.split("/").pop()),
      show: true,
    },
  ];
  return <BreadCrumb crumbs={crumbs} />;
};


const App = ({ path, stateCode, userType, tenants }) => {
  const tenantId = Digit.ULBService.getCurrentTenantId();
  
  const request = {
    url : "/public-service/v1/service",
    headers: {
      "X-Tenant-Id" : tenantId
    },
    method: "GET",
  }
  const {isLoading, data} = Digit.Hooks.useCustomAPIHook(request);
  console.log("dataa",data);

  return (
    <Switch>
      <AppContainer className="ground-container">
        <React.Fragment>
          <SampleBreadCrumbs location={location} />
        </React.Fragment>
        <PrivateRoute path={`${path}/:module/:service/Apply`} component={() => <DigitDemoComponent />} />
        <PrivateRoute path={`${path}/:module/:service/response`} component={() => <Response />} />
        <PrivateRoute path={`${path}/:module/Search`} component={() => <SearchTL />} />
        <PrivateRoute path={`${path}/:module/:service/ViewScreen`} component={() => <DigitDemoViewComponent />} />
        <PrivateRoute path={`${path}/modules`} component={() => <ModulePageComponent />} />
      </AppContainer>
    </Switch>
  );
};

export default App;
