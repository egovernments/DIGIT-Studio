#!/bin/sh

# Run Flyway using uppercase environment variable names
flyway \
  -url="$DB_URL" \
  -table="$SCHEMA_TABLE"  \
  -user="$FLYWAY_USER" \
  -password="$FLYWAY_PASSWORD" \
  -schemas="$FLYWAY_SCHEMAS" \
  -locations="$FLYWAY_LOCATIONS" \
  -baselineOnMigrate=true \
  -outOfOrder=true \
  migrate
