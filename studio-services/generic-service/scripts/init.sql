CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS generic_table (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    value TEXT
    );
