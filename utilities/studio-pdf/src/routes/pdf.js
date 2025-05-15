var express = require("express");
var router = express.Router();
var config = require("../config");
var { search_serviceDetails, create_pdf } = require("../api");
const { asyncMiddleware } = require("../utils/asyncMiddleware");
const { logger } = require("../logger");

function renderError(res, errorMessage, errorCode) {
    if (errorCode === undefined) errorCode = 500;
    res.status(errorCode).send({ errorMessage });
}

router.post(
    "/generate",
    asyncMiddleware(async function (req, res) {
        const tenantId = req.query.tenantId;
        const applicationNumber = req.query.applicationNumber;
        const pdfKey = req.query.pdfKey;
        const serviceCode = req.query.serviceCode
        const requestInfo = req.body;
        

        // Validation
        if (!requestInfo) {
            return renderError(res, "requestinfo cannot be null", 400);
        }
        if (!tenantId) {
            return renderError(res, "tenantId is mandatory to generate the receipt", 400);
        }
        if (!applicationNumber) {
            return renderError(res, "applicationNumber is mandatory to generate the receipt", 400);
        }
        if (!pdfKey) {
            return renderError(res, "pdfKey is mandatory to generate the receipt", 400);
        }

        try {
            // Fetch application details
            let response;
            try {
                response = await search_serviceDetails(tenantId, serviceCode, applicationNumber);
                // logger.info(`response of the application search is ${JSON.stringify(response)}`);
            } catch (ex) {
                logger.error(`Error in the application search for App no ${applicationNumber} of service ${serviceCode}`);
                logger.error(ex);
                return renderError(res, "Failed to query details of the application", 500);
            }

            const application = response.data;
            if (application && application?.Application && application?.Application.length > 0) {
                let pdfResponse;
                try {
                    logger.info(`Sent the request to create a pdf of key ${pdfKey} for application number ${applicationNumber}`);
                    logger.info(`application object  :: ${JSON.stringify(application)}`);
                    pdfResponse = await create_pdf(
                        tenantId,
                        pdfKey,
                        application,
                        requestInfo
                    );

                } catch (ex) {
                    return renderError(res, "Failed to generate PDF for application", 500);
                }

                const filename = `${pdfKey}_${Date.now()}`;
                res.writeHead(200, {
                    "Content-Type": "application/pdf",
                    "Content-Disposition": `attachment; filename=${filename}.pdf`,
                });
                pdfResponse.data.pipe(res);
            } else {
                return renderError(res, "No application found for the given applicationNumber", 404);
            }
        } catch (ex) {
            return renderError(res, "Something went wrong", 500);
        }
    })
);

module.exports = router;
