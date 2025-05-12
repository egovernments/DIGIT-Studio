import React from "react";
import { Card, TextBlock, Button } from "@egovernments/digit-ui-components";
import { transformViewApplication } from "../utils/createUtils";
import { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";

const CheckListCard = (props) => {
    const [filled, setFilled] = useState(false);
    const history = useHistory();

    const style = {
        display: "flex",
        alignItems: "center",
        gap: "1rem",
        margin: "20px"
    };

    const request = {
        url: "/health-service-request/service/v1/_search",
        params: {},
        body: {},
        method: "POST",
        headers: {},
        config: {
            enable: false,
        },
    }
    const mutation = Digit.Hooks.useCustomAPIMutationHook(request);

    const isFilled = async (id, accid) => {

        await mutation.mutate(
            {
                url: '/health-service-request/service/v1/_search',
                method: "POST",
                body: transformViewApplication(id, accid),
                config: {
                    enable: false,
                },
            },
            {
                onSuccess: (res) => {
                    if (res.Services && res.Services.length > 0) {
                        setFilled(true);
                    }
                },
                onError: () => {
                    console.log("Error checking filled status");
                },
            }
        )
    }

    useEffect(() => {
        isFilled(props.item.code, props.accid)
    }, [props.item.code, props.accid]);

    return (
        <Card type="primary" style={style}>
            <TextBlock body={props.item.code} />
            {filled ? (
                <Button label="View Response" onClick={() => history.push({ pathname: `/${window.contextPath}/employee/publicservices/viewresponse` })} />
            ) : (
                <Button label="Fill Checklist" onClick={() => history.push({
                    pathname: `/${window.contextPath}/employee/publicservices/checklist/${props.accid}/${props.item.id}/${props.item.code}`,
                })} />
            )}
        </Card>
    );
};

export default CheckListCard;