CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS egov_service_request (
    id UUID PRIMARY KEY,
    tenantId VARCHAR(250) NOT NULL,
    businessService VARCHAR(250)  NOT NULL,
    module VARCHAR(250)  NOT NULL,
    status VARCHAR(250)  NOT NULL,
    additionalDetail JSONB,
    createdTime TIMESTAMP DEFAULT NOW(),
    lastModifiedTime TIMESTAMP DEFAULT NOW(),
    createdby bigint DEFAULT 1,
    lastmodifiedby bigint DEFAULT 1,
    CONSTRAINT tenantId_businessService_module UNIQUE (tenantId, businessService,module)
);
