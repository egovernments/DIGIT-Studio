import { Link, useHistory } from "react-router-dom";
import _ from "lodash";
import React from "react";

//create functions here based on module name set in mdms(eg->SearchProjectConfig)
//how to call these -> Digit?.Customizations?.[masterName]?.[moduleName]
// these functions will act as middlewares
// var Digit = window.Digit || {};

const businessServiceMap = {};

const inboxModuleNameMap = {};

export const UICustomizations = {
    searchGenericConfig: {
        customValidationCheck: (data) => {
          //checking both to and from date are present
          const { createdFrom, createdTo, field, value } = data;
          if (
            (createdFrom === "" && createdTo !== "") ||
            (createdFrom !== "" && createdTo === "")
          )
            return { type: "warning", label: "ES_COMMON_ENTER_DATE_RANGE" };
    
          if ((field && !value) || (!field && value)) {
            return {
              type: "warning",
              label: "WBH_MDMS_SEARCH_VALIDATION_FIELD_VALUE_PAIR",
            };
          }
    
          return false;
        },
        preProcess: (data, additionalDetails) => {
          const tenantId = Digit.ULBService.getCurrentTenantId();
          data.body.SearchCriteria.tenantId = tenantId;
        //   const filters = {};
          const custom = data.body.SearchCriteria.custom;
          data.headers = {"X-Tenant-Id":tenantId}
          data.method = "GET"
          if(custom?.applicationNumber) data.params.applicationNumber = custom?.applicationNumber;
          if(custom?.status) data.params.status = custom?.status;
          if(custom?.todate) data.params.todate = Digit.Utils.date.convertDateToEpoch(custom?.todate);
          if(custom?.fromdate) data.params.todate = Digit.Utils.date.convertDateToEpoch(custom?.fromdate);
          if(data?.state?.searchForm?.businessService) data.url = `${data.url}/${data?.state?.searchForm?.businessService?.serviceCode}`
          delete data.body.SearchCriteria.custom;
        //   const { field, value, isActive } = custom || {};
        //   filters[field?.code] = value;
        //   if (isActive) {
        //     if (isActive.value === "all") delete data.body.MdmsCriteria.isActive;
        //     else data.body.MdmsCriteria.isActive = isActive?.value;
        //   } else {
        //     delete data.body.MdmsCriteria.isActive;
        //   }
        //   data.body.MdmsCriteria.filters = filters;
        //   // data.body.MdmsCriteria.limit = 100
        //   data.body.MdmsCriteria.limit = data.state.tableForm.limit;
        //   data.body.MdmsCriteria.offset = data.state.tableForm.offset;
        //   data.body.MdmsCriteria.schemaCode =
        //     // additionalDetails?.currentSchemaCode
        //     "ACCESSCONTROL-ACTIONS-TEST.actions-test";
        //   delete data.body.MdmsCriteria.custom;
          return data;
        },
        additionalCustomizations: (row, key, column, value, t, searchResult) => {
          switch (key) {
            case `${row?.module.toUpperCase()}_APPLICATION_NUMBER`:
              return (
                <span className="link">
            <Link to={`/${window.contextPath}/employee/publicservices/${row?.module}/${row?.businessService}/ViewScreen?applicationNumber=${row?.applicationNumber}&serviceCode=${row?.serviceCode}`}>
            {String(value ? value : t("ES_COMMON_NA"))}
            </Link>
          </span>
              );
            default:
              return t("ES_COMMON_NA");
          }
        },
        MobileDetailsOnClick: (row, tenantId) => {
          let link;
          Object.keys(row).map((key) => {
            if (key === "MASTERS_WAGESEEKER_ID")
              link = `/${window.contextPath}/employee/masters/view-wageseeker?tenantId=${tenantId}&wageseekerId=${row[key]}`;
          });
          return link;
        },
        additionalValidations: (type, data, keys) => {
          if (type === "date") {
            return data[keys.start] && data[keys.end]
              ? () =>
                  new Date(data[keys.start]).getTime() <=
                  new Date(data[keys.end]).getTime()
              : true;
          }
        },
        selectionHandler: (event) => {
        }, // selectionHandler : Is used to handle row selections. gets on object which containes 3 key value pairs:  allSelected(whether all rows are selected or not), selectedCount (no, of rows selected),selectedRows( an array of selected rows)
        actionSelectHandler: (index, label, selectedRows) => {
        }, // actionSelectHandler : Is used to handle onClick functions of table action button on row selections, gets index,label and selectedRows as props
        footerActionHandler: (index, event) => {
        }, // footerActionHandler : Is used to handle onclick functions of footer action buttons, gets index and event as props
        linkColumnHandler: (row) => {
          const url = `/${window.contextPath}/employee/microplan/view-main?tenantId=${row?.tenantId}&uniqueIdentifier=${row?.uniqueIdentifier}`;
          window.location.href = url;
        }, 
      },
};