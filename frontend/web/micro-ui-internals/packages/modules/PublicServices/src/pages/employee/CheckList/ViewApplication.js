import React from "react";
import { Card, TextBlock, Button } from "@egovernments/digit-ui-components";
import { transformViewApplication } from "../../../utils/createUtils";
import { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";

const ViewApplication = () => {
    
    const history = useHistory();
    const code = [
        "SMC BHAVYA.TRAINING_SUPERVISION.TEAM_SUPERVISOR",
        "LLIN-mz_april_2025.TRAINING_SUPERVISION.PROVINCIAL_SUPERVISOR"
    ];
    const accountID = "";
    const [filledStatus, setFilledStatus] = useState({});

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
            enable: true,
        },
    }
    const mutation = Digit.Hooks.useCustomAPIMutationHook(request);

    const isFilled = async (id, accid) => {

        console.log("isfilled");
        await mutation.mutate(
            {
                url: '/health-service-request/service/v1/_search',
                method: "POST",
                body: transformViewApplication(id, accid),
                config: {
                    enable: true,
                },
            },
            {
                onSuccess: (res) => {
                    console.log(res, "response");
                    const filled = Array.isArray(res?.services) && res.services.length > 0;
                    setFilledStatus((prev) => ({ ...prev, [id]: filled }));
                },
                onError: () => {
                    setFilledStatus((prev) => ({ ...prev, [id]: false }));
                },
            }
        )
    }

    useEffect(() => {
        code.forEach((c) => isFilled(c, accountID));
    }, [accountID]);

    return (
        <React.Fragment>
            {code.map((item, index) => (
                <Card type="primary" key={index} style={style}>
                    <TextBlock body={item} />
                    {false ? (
                        <Button label="View Response" onClick={() => history.push({ pathname: `/${window.contextPath}/employee/publicservices/viewresponse` })} />
                    ) : (
                        <Button label="Fill Checklist" onClick={() => history.push({ pathname: `/${window.contextPath}/employee/publicservices/checklist` })} />
                    )}
                </Card>
            ))}
        </React.Fragment>
    );
};

export default ViewApplication;