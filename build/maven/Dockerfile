# syntax=docker/dockerfile:1

# Use official Golang image to build
FROM golang:1.21-alpine AS build

# Accept working directory and TYPE from build args
ARG WORK_DIR
ARG SERVICE_TYPE=app  # default is 'app', can be overridden to 'db'

ENV WORK_DIR=${WORK_DIR}
ENV SERVICE_TYPE=${SERVICE_TYPE}

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (if exists)
COPY ${WORK_DIR}/go.mod ${WORK_DIR}/go.sum ./

# Download dependencies (only if go.mod exists)
RUN if [ -f go.mod ]; then go mod download; fi

# Copy full source code
COPY ${WORK_DIR} .

# Build based on type
RUN if [ "$SERVICE_TYPE" = "db" ]; then \
      go build -o /app/migration scripts/migration.go; \
    else \
      go build -o /app/main main.go; \
    fi

# Final minimal image
FROM alpine:latest

# Install tini to properly handle PID 1 (optional but good practice)
RUN apk add --no-cache tini

# Set working directory
WORKDIR /opt/egov

# Copy compiled binaries
COPY --from=build /app/main /opt/egov/main
COPY --from=build /app/migration /opt/egov/migration

# Copy SQL migrations if any
COPY ${WORK_DIR}/scripts/migration /opt/egov/migration-scripts/

# Expose default app port (optional)
EXPOSE 8080

# Entry point based on type
ENTRYPOINT ["/sbin/tini", "--"]

CMD if [ "$SERVICE_TYPE" = "db" ]; then \
      echo "Running database migration..." && ./migration; \
    else \
      echo "Starting Go application..." && ./main; \
    fi
