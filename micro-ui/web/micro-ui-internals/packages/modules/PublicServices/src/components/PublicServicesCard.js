import { HRIcon, EmployeeModuleCard, AttendanceIcon, PropertyHouse } from "@egovernments/digit-ui-react-components";
import React from "react";
import { useTranslation } from "react-i18next";

const PublicServicesCard = () => {
 
  const { t } = useTranslation();

  const propsForModuleCard = {
    Icon: "BeenHere",
    moduleName: t("Digit Studio"),
    kpis: [

    ],
    links: [
      {
        label: t("Services Apply (TL)"),
        link: `/${window?.contextPath}/employee/publicservices/digit-demo-apply/tl`,
      },
    ],
  };

  return <EmployeeModuleCard {...propsForModuleCard} />;
};

export default PublicServicesCard;