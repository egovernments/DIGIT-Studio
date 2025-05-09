import React from "react";
import PropTypes from "prop-types";
import { useTranslation } from "react-i18next";
import LabelFieldPair from "../atoms/LabelFieldPair";
import CardLabel from "../atoms/CardLabel";
import CardLabelError from "../atoms/CardLabelError";
import CitizenInfoLabel from "../atoms/CitizenInfoLabel";
import HeaderComponent from "../atoms/HeaderComponent";
import MultiUploadWrapper from "../molecules/MultiUploadWrapper";
import TextInput from "../atoms/TextInput";
import { getRegex } from "../utils/uploadFileComposerUtils";
import { Loader } from "@egovernments/digit-ui-react-components";
import { useParams } from "react-router-dom/cjs/react-router-dom.min";

const UploadAndDownloadDocumentHandler = ({
  schemaCode = "DigitStudio.DocumentConfig",
  config,
  Controller,
  control,
  register,
  formData,
  errors,
  localePrefix,
  customClass,
  action="APPLY",
  flow
}) => {
  const { t } = useTranslation();
  const tenantId = Digit?.ULBService?.getStateId();
  const { module, service} = useParams();
  let moduleName = `${module?.toLowerCase()}.${service?.toLowerCase()}`;

  const requestCriteria = {
    url: "/egov-mdms-service/v1/_search",
    body: {
      MdmsCriteria: {
        "tenantId": tenantId,
        "moduleDetails": [
            {
                "moduleName": "DigitStudio",
                "masterDetails": [
                    {
                        "name": "DocumentConfig2"
                    }
                ]
            }
        ]
    },
    },
    changeQueryName: "documentConfig",
  };

  const { isLoading, data } = Digit.Hooks.useCustomAPIHook(requestCriteria);

let docData = data ? data?.MdmsRes?.DigitStudio?.DocumentConfig2?.filter((ob) => ob?.module.toLowerCase() === moduleName)?.[0]?.actions : [];

const docConfig = docData?.filter((item) => item?.action === "APPLY")?.[0];
  if (!docConfig && flow !== "WORKFLOW") return null;
  if(isLoading) return <Loader />;
  return (
    <React.Fragment>
      {/* <HeaderComponent styles={{ fontSize: "24px", marginTop: "40px" }}>
        {t("WORKS_RELEVANT_DOCUMENTS")}
      </HeaderComponent> */}

      {/* {docConfig?.bannerLabel && (
        <CitizenInfoLabel
          info={t("ES_COMMON_INFO")}
          text={t(docConfig?.bannerLabel)}
          className="digit-doc-banner"
        />
      )} */}
      {flow === "WORKFLOW" && 
            <Controller
              name={`${config?.populators?.name}`}
              control={control}
              rules={{ required: false }}
              render={({ onChange, ref, value = [] }) => {
                function getFileStoreData(filesData) {
                  const numberOfFiles = filesData.length;
                  let finalDocumentData = [];
                  if (numberOfFiles > 0) {
                    filesData.forEach((value) => {
                      finalDocumentData.push({
                        fileName: value?.[0],
                        fileStoreId: value?.[1]?.fileStoreId?.fileStoreId,
                        documentType: value?.[1]?.file?.type,
                      });
                    });
                  }
                  onChange(numberOfFiles > 0 ? filesData : []);
                }
                return (
                  <MultiUploadWrapper
                    t={t}
                    module="works"
                    tenantId={Digit.ULBService.getCurrentTenantId()}
                    getFormState={getFileStoreData}
                    showHintBelow={config?.populators?.showHintBelow ? true : false}
                    setuploadedstate={value}
                    allowedFileTypesRegex={getRegex(config?.populators?.allowedFileTypes)}
                    allowedMaxSizeInMB={config?.populators?.maxSizeInMB}
                    hintText={t(config?.populators?.hintText)}
                    maxFilesAllowed={config?.populators?.maxFilesAllowed}
                    extraStyleName={{ padding: "0.5rem" }}
                    //customClass={populators?.customClass}
                  />
                );
              }}
            />
      }
      {flow !== "WORKFLOW" && docConfig?.documents?.map((item, index) => {
        if (!item?.active) return null;
        return (
          <LabelFieldPair key={index} style={{ alignItems: item?.showTextInput ? "flex-start" : "center" }}>
            {item.code && (
              <div style={{ display: "flex", flexDirection: "column" }}>
                <CardLabel className="bolder" style={{ marginTop: item?.showTextInput ? "10px" : "" }}>
                  {t(`${localePrefix}_${item?.code}`)} {item?.isMandatory ? " * " : null}
                </CardLabel>

                {(item?.templatePDFKey || item?.templatedownloadURL) && (
           <div style={{ display: "flex", alignItems: "center", gap: "1rem", marginBottom: "1rem", width: "100%" }}>
         
           <div className={`digit-upload-wrapper ${customClass || ""}`} style={{ flex: 1, padding: "1rem", border: "1px solid #D6D5D4", borderRadius: "8px", backgroundColor: "#FAFAFA", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
             <div style={{ fontSize: "14px", color: "#1A1A1A" }}>{t(item?.documentType)}</div>
             <button
               type="button"
               onClick={async () => {
                 try {
                   const state = tenantId;
                   if (item?.templatedownloadURL) {
                     window.open(item.templatedownloadURL, "_blank");
                   } else if (item?.templatePDFKey) {
                     const dummyPayload = { sample: "value" };
                     const response = await Digit.PaymentService.generatePdf(state, dummyPayload, item.templatePDFKey);
                     const fileStore = await Digit.PaymentService.printReciept(state, {
                       fileStoreIds: response.filestoreIds[0],
                     });
                     const fileUrl = fileStore?.[response.filestoreIds[0]];
                     if (fileUrl) {
                       window.open(fileUrl, "_blank");
                     }
                   }
                 } catch (err) {
                   console.error("Template download error", err);
                 }
               }}
               style={{
                 fontSize: "14px",
                 padding: "6px 12px",
                 backgroundColor: "#007AFF",
                 color: "#fff",
                 border: "none",
                 borderRadius: "4px",
                 cursor: "pointer",
               }}
             >
               {t("DOWNLOAD_TEMPLATE")}
             </button>
           </div>
         </div>
         
          
)}
              </div>
            )}

            <div className="digit-field">
              {item?.showTextInput && (
                <TextInput
                  style={{ marginBottom: "16px" }}
                  name={`${config?.name}.${item?.name}_name`}
                  placeholder={t("ES_COMMON_ENTER_NAME")}
                  inputRef={register({ minLength: 2 })}
                />
              )}

              {!(item?.templatePDFKey || item?.templatedownloadURL) && <div style={{ marginBottom: "24px" }}>
                <Controller
                  render={({ value = [], onChange }) => {
                    function getFileStoreData(filesData) {
                      let finalDocumentData = [];
                      filesData.forEach((value) => {
                        finalDocumentData.push({
                          fileName: value?.[0],
                          fileStoreId: value?.[1]?.fileStoreId?.fileStoreId,
                          documentType: value?.[1]?.file?.type,
                        });
                      });
                      onChange(finalDocumentData.length ? filesData : []);
                    }

                    return (
                      <MultiUploadWrapper
                        t={t}
                        module="DigitStudio"
                        getFormState={getFileStoreData}
                        setuploadedstate={value}
                        showHintBelow={Boolean(item?.hintText)}
                        hintText={item?.hintText}
                        allowedFileTypesRegex={getRegex(item?.allowedFileTypes)}
                        allowedMaxSizeInMB={item?.maxSizeInMB || 5}
                        maxFilesAllowed={item?.maxFilesAllowed || 1}
                        customErrorMsg={item?.customErrorMsg}
                        customClass={customClass}
                        tenantId={Digit.ULBService.getCurrentTenantId()}
                      />
                    );
                  }}
                  rules={{
                    validate: (value) => !(item?.isMandatory && (!value || value.length === 0)),
                  }}
                  defaultValue={formData?.[item?.name]}
                  name={`${config?.name}.${item?.name}`}
                  control={control}
                />
                {errors?.[`${config?.name}`]?.[`${item?.name}`] && (
                  <CardLabelError style={{ fontSize: "12px" }}>
                    {t(config?.error)}
                  </CardLabelError>
                )}
              </div>}
            </div>
          </LabelFieldPair>
        );
      })}
    </React.Fragment>
  );
};

UploadAndDownloadDocumentHandler.propTypes = {
  schemaCode: PropTypes.string.isRequired,
  config: PropTypes.object.isRequired,
  Controller: PropTypes.func.isRequired,
  control: PropTypes.object.isRequired,
  register: PropTypes.func.isRequired,
  formData: PropTypes.object.isRequired,
  errors: PropTypes.object.isRequired,
  localePrefix: PropTypes.string.isRequired,
  customClass: PropTypes.string,
};

export default UploadAndDownloadDocumentHandler;
