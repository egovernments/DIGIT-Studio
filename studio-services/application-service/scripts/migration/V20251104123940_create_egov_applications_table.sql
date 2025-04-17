-- Create table: eg_applications
CREATE TABLE IF NOT EXISTS eg_applications (
    id UUID PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    module VARCHAR(255) NOT NULL,
    business_service VARCHAR(255) NOT NULL,
    status VARCHAR(255),
    channel VARCHAR(255),
    application_number VARCHAR(255),
    workflow_status VARCHAR(255),
    service_details JSONB,
    additional_details JSONB,
    address JSONB,
    workflow JSONB,
    createdby VARCHAR(255),
    last_modifiedby VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_tenant_module_businessservice UNIQUE (tenant_id, module, business_service)
    );

-- Add unique constraint on tenant_id + module + business_service
/*ALTER TABLE eg_applications
    ADD CONSTRAINT uniq_tenant_module_businessservice UNIQUE (tenant_id, module, business_service);*/

--------------------------------------------------------------------------------------

-- Create table: eg_reference
CREATE TABLE IF NOT EXISTS eg_reference (
    id UUID PRIMARY KEY,
    reference_type VARCHAR(255),
    module VARCHAR(255),
    tenant_id VARCHAR(255),
    reference_no VARCHAR(255),
    active BOOLEAN DEFAULT FALSE,
    application_id UUID,
    CONSTRAINT fk_reference_application FOREIGN KEY (application_id) REFERENCES eg_applications(id) ON DELETE CASCADE
    );

--------------------------------------------------------------------------------------

-- Create table: eg_applicant
CREATE TABLE IF NOT EXISTS eg_applicant (
    id UUID PRIMARY KEY,
    type VARCHAR(255),
    application_id UUID,
    user_id VARCHAR(255),
    name VARCHAR(255),
    mobile_number BIGINT,
    email_id VARCHAR(255),
    prefix VARCHAR(50),
    active BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_applicant_application FOREIGN KEY (application_id) REFERENCES eg_applications(id) ON DELETE CASCADE
    );
---------------------------------------

CREATE TABLE IF NOT EXISTS eg_application_sequence (
                                                       id SERIAL PRIMARY KEY,
                                                       tenant_id VARCHAR(255) NOT NULL,
    module VARCHAR(255) NOT NULL,
    business_service VARCHAR(255) NOT NULL,
    last_number INT DEFAULT 0,
    UNIQUE (tenant_id, module, business_service)
    );