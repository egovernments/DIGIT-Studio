import { Header, ActionBar, SubmitBar, Menu, Card, Loader, ViewComposer, MultiLink } from "@egovernments/digit-ui-react-components";
import React, { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
import { generateViewConfigFromResponse } from "../../../utils";
import WorkflowActions from "../../../components/WorkflowActions";
import ViewCheckListCards from "../CheckList/viewCheckListCards";
import { useWorkflowDetailsWorks } from "../../../utils";

const DigitDemoViewComponent = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const tenantId = Digit.ULBService.getCurrentTenantId();

  const userInfo = Digit.UserService.getUser();
  const userRoles = userInfo?.info?.roles?.map((roleData) => roleData?.code);
  const [current, setCurrent] = useState(Date.now());
  const queryStrings = Digit.Hooks.useQueryParams();

  const request = {
    url : `/public-service/v1/application/${queryStrings?.serviceCode|| "SVC-DEV-TRADELICENSE-NEWTL-04"}`,
    headers: {
      "X-Tenant-Id" : tenantId
    },
    method: "GET",
    params: {
      "applicationNumber": queryStrings?.applicationNumber,
    }
  }
  const {isLoading, data} = Digit.Hooks.useCustomAPIHook(request);

  let response =  data ? data?.Application?.[0] : {};
  let config = generateViewConfigFromResponse(response,t);

  let {data :workflowDetails, isLoading: workflowLoading} = useWorkflowDetailsWorks(
    {
      tenantId: tenantId,
      id: queryStrings?.applicationNumber,
      moduleCode: response?.businessService,
      config: {
        enabled: response ? true : false,
        cacheTime: 0
      }
    }
  );
  let checkListCodes = workflowDetails ? [`${response?.businessService}.${workflowDetails?.processInstances[0].state?.state}`] : [];
  if (isLoading || workflowLoading) {
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
      <ViewCheckListCards applicationId={data?.Application?.[0]?.id} checkListCodes={checkListCodes} />
      <WorkflowActions
          forcedActionPrefix={`WF_${response?.businessService}_ACTION`}
          businessService={response?.businessService}
          applicationNo={response?.applicationNumber}
          tenantId={tenantId}
          applicationDetails={response}
          url={`/public-service/v1/application/SVC-DEV-TRADELICENSE-NEWTL-04/${response?.id}`}
          //setStateChanged={setStateChanged}
          moduleCode={response?.module}
          //editApplicationNumber={""}
          // WorflowValidation={(setShowModal) => {
          //   try {
          //     let validationFlag = false;
          //     for (const validation of mbValidationMr?.musterRollValidation) {
          //       if (validation?.type === "error") {
          //         validationFlag = true;
          //         setShowToast({ type: "error", label: t(validation?.message) });
          //         break;
          //       } else if (validation?.type === "warn") {
          //         validationFlag = true;
          //         setShowPopUp({ setShowWfModal: setShowModal, label: t(validation?.message) });
          //         break;
          //       }
          //     }
          //     if (!validationFlag) setShowModal(true);
          //   } catch (error) {
          //     showToast(error.message);
          //   }
          // }}
          // editCallback={() => {
          //   setModify(true);
          //   setshowEditTitle(true);
          //   setSaveAttendanceState((prevState) => {
          //     return {
          //       ...prevState,
          //       displaySave: true,
          //       updatePayload: data?.applicationData?.individualEntries?.map((row) => {
          //         return {
          //           totalAttendance: row?.modifiedTotalAttendance || row?.actualTotalAttendance,
          //           id: row?.id,
          //         };
          //       }),
          //     };
          //   });
          // }}
        />
      {/* <ActionBar>
        {displayMenu ? <Menu localeKeyPrefix={"WORKS"} options={actionULB} optionKey={"name"} t={t} onSelect={onActionSelect} /> : null}
        <SubmitBar label={t("ACTIONS")} onSubmit={() => setDisplayMenu(!displayMenu)} />
      </ActionBar> */}
    </React.Fragment>
  );
};

export default DigitDemoViewComponent;
