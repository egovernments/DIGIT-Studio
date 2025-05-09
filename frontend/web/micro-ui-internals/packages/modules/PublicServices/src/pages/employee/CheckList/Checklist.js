import React, { useReducer } from "react";
import { useTranslation } from "react-i18next";
import { FormComposerV2 } from "@egovernments/digit-ui-components";
import { CheckListConfig } from "../../../configs/checklistconfig";

const formReducer = (state, action) => {
    switch (action.type) {
        case 'UPDATE_FORM':
            return {
                ...state,
                formData: action.payload
            };
        default:
            return state;
    }
};

const CheckList = () => {
    //const { t } = useTranslation();
    const [state, dispatch] = useReducer(formReducer, {
        formData: {}
    });

    const config = CheckListConfig(state.formData);
    const onSubmit = async (data) => {
        console.log(data, "data");
    };

    return (
        <div>
            <FormComposerV2
                //label={t("SUBMIT")}
                config={config}
                defaultValues={state.formData}
                onFormValueChange={(setValue, formData) => {
                    console.log(formData, "formdata");
                    if (JSON.stringify(formData) !== JSON.stringify(state.formData)) {
                        console.log("asf");
                        dispatch({
                            type: 'UPDATE_FORM',
                            payload: formData
                        });
                    }
                }}
                onSubmit={onSubmit}
                //fieldStyle={{ marginRight: 2 }}
            />
        </div>
    );
};

export default CheckList;
