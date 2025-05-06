import { HRIcon, EmployeeModuleCard, AttendanceIcon, PropertyHouse } from "@egovernments/digit-ui-react-components";
import React from "react";
import { useTranslation } from "react-i18next";

const PublicServicesCard = () => {
 
  const { t } = useTranslation();

  const propsForModuleCard = {
    Icon: "BeenHere",
    moduleName: t("DIGIT_STUDIO"),
    kpis: [

    ],
    links: [
      {
        label: t("DIGIT_STUDIO_APPLY"),
        link: `/${window?.contextPath}/employee/publicservices/modules?selectedPath=Apply`,
      },
      // {
      //   label: t("Services Search"),
      //   link: `/${window?.contextPath}/employee/publicservices/modules?selectedpath=Search`,
      // },
      {
        label: t("Services Apply (PGR)"),
        link: `/${window?.contextPath}/employee/publicservices/pgr/Newpgr/Apply`,
      },
      {
        label: t("Services Search (TL)"),
        link: `/${window?.contextPath}/employee/publicservices/tl/Search`,
      },
      {
        label: t("Services Search (PGR)"),
        link: `/${window?.contextPath}/employee/publicservices/pgr/Search`,
      },
      {
        label: t("Services Inbox Search (TL)"),
        link: `/${window?.contextPath}/employee/publicservices/tl/Inbox`,
      },
      {
        label: t("Services Inbox Search (PGR)"),
        link: `/${window?.contextPath}/employee/publicservices/pgr/Inbox`,
      },
    ],
  };

  return <EmployeeModuleCard {...propsForModuleCard} />;
};

export default PublicServicesCard;