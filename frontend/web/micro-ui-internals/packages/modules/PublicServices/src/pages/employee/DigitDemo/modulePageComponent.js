import React from "react";
import { Card, Button, HeaderComponent, CardText, Loader, SubmitBar } from "@egovernments/digit-ui-components";
import { useTranslation } from "react-i18next";
import { useHistory } from "react-router-dom";
import { transformResponseforModulePage } from "../../../utils";
import { Link } from "react-router-dom/cjs/react-router-dom.min";

const modulePageComponent = ({}) => {
  const { t } = useTranslation();
  const history = useHistory();

  const tenantId = Digit.ULBService.getCurrentTenantId();
  const queryStrings = Digit.Hooks.useQueryParams();

  const request = {
    url : "/public-service/v1/service",
    headers: {
      "X-Tenant-Id" : tenantId
    },
    method: "GET",
  }
  const {isLoading, data} = Digit.Hooks.useCustomAPIHook(request);
  console.log(data);

  let detailsConfig = data ? transformResponseforModulePage(data?.Services) : [];
console.log(detailsConfig);
  if (isLoading) {
    return <Loader />;
  }

  return (
    <div className="products-container">
      {/* Header Section */}
      <HeaderComponent className="products-title">{t("DIGIT_STUDIO_HEADER")}</HeaderComponent>
      <CardText className="products-description">
        {t("DIGIT_STUDIO_HEADER_DESCRIPTION")}
      </CardText>

      {/* Product Cards Section */}
      <div className="products-list">
        {detailsConfig?.map((product, index) => (
          <Card key={index} className="product-card">
            <div className="product-header">
              <HeaderComponent className="product-title">{t(product.heading)}</HeaderComponent>
            </div>
            <CardText className="product-description">{t(product?.cardDescription)}</CardText>
            {queryStrings?.selectedPath === "Apply" && product?.businessServices.map((bs) => (
             <Link className="link" to={`/${window.contextPath}/employee/publicservices/${product.module}/${bs.businessService}/Apply?serviceCode=${bs?.serviceCode}`}>
              {bs.businessService}
        </Link>
            ))
            }
            <Link className="link" to={`/${window.contextPath}/employee/publicservices/${product.module}/Search`}>
              Search
            </Link>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default modulePageComponent;