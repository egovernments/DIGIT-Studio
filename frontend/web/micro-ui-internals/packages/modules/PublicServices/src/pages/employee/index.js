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
import InboxService from "./DigitDemo/InboxService";
import ViewApplication from "./CheckList/viewApplication";
import CheckList from "./CheckList/Checklist";

const SampleBreadCrumbs = ({ location }) => {
  const { t } = useTranslation();

  const crumbs = [
    {
      internalLink: `/${window?.contextPath}/employee`,
      content: t("HOME"),
      show: true,
    },
    {
      content: t(location.pathname.split("/").pop().toUpperCase()),
      show: true,
    }
  ];
  return <BreadCrumb crumbs={crumbs} />;
};


const App = ({ path, stateCode, userType, tenants }) => {
  const tenantId = Digit.ULBService.getCurrentTenantId();

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
        <PrivateRoute path={`${path}/:module/inbox`} component={() => <InboxService />} />
        <PrivateRoute path={`${path}/viewapp`} component={() => <ViewApplication />} />
        <PrivateRoute path={`${path}/checklist`} component={() => <CheckList />} />
      </AppContainer>
    </Switch>
  );
};

export default App;
