import { AppContainer, PrivateRoute } from "@egovernments/digit-ui-react-components";
import { BreadCrumb } from "@egovernments/digit-ui-components";
import React from "react";
import { useTranslation } from "react-i18next";
import { Switch } from "react-router-dom";
// import Inbox from "./SampleInbox";
import DigitDemoComponent from "./DigitDemo/digitDemoComponent";
import SearchTL from "./DigitDemo/searchTL";

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

  const tenantId=`dev`;
  const request = {
    url : "/public-service/v1/service",
    method: "GET",
    headers: {
      "X-Tenant-Id" : tenantId ,
    },
  }
  const {isLoading, data} = Digit.Hooks.useCustomAPIHook(request);
  console.log("dattaa",data);

  return (
    <Switch>
      <AppContainer className="ground-container">
        <React.Fragment>
          <SampleBreadCrumbs location={location} />
        </React.Fragment>
        <PrivateRoute path={`${path}/:module/Apply`} component={() => <DigitDemoComponent />} />
        <PrivateRoute path={`${path}/tl/Search`} component={() => <SearchTL />} />
      </AppContainer>
    </Switch>
  );
};

export default App;
