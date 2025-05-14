// import Dropdown from '../../Dropdown';
// import { Loader } from '../../Loader';
// import React, { useState } from 'react'


const configModal = (
    t,
    action,
    approvers,
    businessService,
    moduleCode,
    documentConfig = []
) => {
    
    const {action:actionString} = action
    const prefix = `${moduleCode.toUpperCase()}_${businessService.toUpperCase()}`;

    let currentModule = `${moduleCode.toLowerCase()}.${businessService.toLowerCase()}`;
    let docData = documentConfig?.filter((ob) => ob?.module.toLowerCase() === currentModule)?.[0]?.actions;
  
    let docConfig = docData?.filter((item) => item?.action === actionString)?.[0];
    if(docConfig === undefined) docConfig = docData?.filter((item) => item?.action === "DEFAULT")?.[0]
//field can have (comments,assignee,upload)
    const fetchIsMandatory = (field) => {
            return docConfig?.[field]?.isMandatory ? docConfig?.[field]?.isMandatory : false
    }
    const fetchIsShow = (field) => {
           return docConfig?.[field]?.show ? docConfig?.[field]?.show : false
    }

    const documentFields = (docConfig?.documents || []).map((doc, index) => ({
        type: "documentUploadAndDownload",
        label: t(`${doc.code}`),
        localePrefix: prefix,
        populators: {
            name: `document.${doc.name}` || `document_${index}`,
            allowedMaxSizeInMB: doc.maxSizeInMB || 5,
            maxFilesAllowed: doc.maxFilesAllowed || 2,
            allowedFileTypes: doc.allowedFileTypes,
            hintText: t(doc.hintText || "COMMON_DOC_UPLOAD_HINT"),
            showHintBelow: true,
            customClass: "upload-margin-bottom",
            errorMessage: t(doc.errorMessage || "COMMON_FILE_UPLOAD_CUSTOM_ERROR_MSG"),
            hideInForm: false,
            action: actionString,
            flow: "WORKFLOW"
        }
    }));

    return {
        label: {
            heading: Digit.Utils.locale.getTransformedLocale(`WF_MODAL_HEADER_${businessService}_${action.action}`),
            submit: Digit.Utils.locale.getTransformedLocale(`WF_MODAL_SUBMIT_${businessService}_${action.action}`),
            cancel: "WF_MODAL_CANCEL",
        },
        form: [
            {
                body: [
                    {
                        label: " ",
                        type: "checkbox",
                        disable: false,
                        isMandatory:false,
                        populators: {
                            name: "acceptTerms",
                            title: "MUSTOR_APPROVAL_CHECKBOX",
                            isMandatory: false,
                            labelStyles: {},
                            customLabelMarkup: true,
                            hideInForm: !fetchIsShow("acceptTerms")
                        }
                    },
                    {
                        label: t("WF_MODAL_APPROVER"),
                        type: "dropdown",
                        isMandatory: fetchIsMandatory("assignee"),
                        disable: false,
                        key:"assignees",
                        populators: {
                            name: "assignee",
                            optionsKey: "nameOfEmp",
                            options: approvers,
                            hideInForm: !fetchIsShow("assignee"),
                            "optionsCustomStyle": {
                                "top": "2.3rem"
                              }
                        },
                    },
                    {
                        label: t("WF_MODAL_COMMENTS"),
                        type: "textarea",
                        isMandatory: fetchIsMandatory("comments"),
                        populators: {
                            name: "comments",
                            hideInForm:!fetchIsShow("comments"),
                            validation:{
                                maxLength:{
                                    value:1024,
                                    message:t("COMMON_COMMENT_LENGTH_EXCEEDED_1024")
                                }
                            }
                        },
                    },
                    ...documentFields,
                ]
            }
        ]
    }
}

export default configModal