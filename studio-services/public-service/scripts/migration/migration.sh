#!/bin/sh

# Load environment variables from .env file
set -a
. ./env/.env
set +a

# Run Flyway using uppercase environment variable names
flyway \
  -url="$FLYWAY_URL" \
  -user="$FLYWAY_USER" \
  -password="$FLYWAY_PASSWORD" \
  -schemas="$FLYWAY_SCHEMAS" \
  -locations="$FLYWAY_LOCATIONS" \
  -baselineOnMigrate=true \
  migrate
