import React from "react";
import { Card, TextBlock, Button } from "@egovernments/digit-ui-components";

const ViewApplication = () => {
    const code = [
        "SMC BHAVYA.TRAINING_SUPERVISION.TEAM_SUPERVISOR",
        "LLIN-mz_april_2025.TRAINING_SUPERVISION.PROVINCIAL_SUPERVISOR"
    ];
    const accountID = "";

    const style = {
        display: "flex",
        alignItems: "center",
        gap: "1rem",
        margin: "20px"
    };

    

    return (
        <React.Fragment>
            {code.map((item, index) => (
                <Card type="primary" key={index} style={style}>
                    <TextBlock body={item} />
                    {true ? (
                        <Button label="View Response" />
                    ) : (
                        <Button label="Fill Checklist" />
                    )}
                </Card>
            ))}
        </React.Fragment>
    );
};

export default ViewApplication;
