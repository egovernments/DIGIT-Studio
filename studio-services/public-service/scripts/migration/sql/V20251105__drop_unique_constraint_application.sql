ALTER TABLE application
DROP CONSTRAINT IF EXISTS uniq_tenant_module_businessservice;

-- Add the new unique constraint
ALTER TABLE application
    ADD CONSTRAINT uniq_service_code_application_number UNIQUE (service_code, application_number);
