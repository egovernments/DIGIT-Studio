# Public Service

The Public Service is a microservice that manages public services and applications within the DIGIT ecosystem. It provides APIs for creating, updating, and searching services and applications, with workflow integration.

## Features

- Service management (create/update/search)
- Application management with workflow integration
- Integration with Individual, Billing, and MDMS services
- Kafka-based asynchronous processing
- PostgreSQL data persistence with Flyway migrations

## Architecture

The service follows a clean architecture pattern with distinct layers:

- **Controllers**: Handle HTTP requests and responses
- **Services**: Implement business logic
- **Repositories**: Handle data access and persistence
- **Models**: Define domain entities and data transfer objects
- **Utils**: Provide common utilities for logging and error handling

## Setup and Installation

### Prerequisites

- Go 1.24+
- PostgreSQL 15
- Kafka

### Environment Configuration

Configure the service using environment variables in a `.env` file:

Server Configuration
SERVER_PORT=8080

Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SCHEMA=public

External Service Endpoints
WORKFLOW_HOST=http://localhost:8081/
WORKFLOW_TRANSITION_PATH=egov-workflow-v2/egov-wf/process/_transition
WORKFLOW_SEARCH_PATH=egov-workflow-v2/egov-wf/process/_search

... more environment variables

## API Endpoints

### Service Management

#### Create Service
- **Endpoint**: POST `/public-service/v1/service`
- **Required Headers**: `X-Tenant-Id`

#### Update Service
- **Endpoint**: PUT `/public-service/v1/service/{serviceCode}`
- **Required Headers**: `X-Tenant-Id`

#### Search Services
- **Endpoint**: GET `/public-service/v1/service`
- **Required Headers**: `X-Tenant-Id`
- **Query Parameters**: `module`, `businessService`, `serviceCode`

### Application Management

#### Create Application
- **Endpoint**: POST `/public-service/v1/application/{serviceCode}`
- **Required Headers**: `X-Tenant-Id`

#### Update Application
- **Endpoint**: PUT `/public-service/v1/application/{serviceCode}/{applicationId}`
- **Required Headers**: `X-Tenant-Id`

#### Search Applications
- **Endpoint**: GET `/public-service/v1/application/{serviceCode}`
- **Required Headers**: `X-Tenant-Id`
- **Query Parameters**: `module`, `businessService`, `status`, `applicationNumber`, `ids`



## 📄 Documentation

- 📐 [Design Document](https://docs.google.com/document/d/13LR7TQMsIg0nD5-Wdl4kj1r3kYjzLyKD0FVzvJkkR3s/edit?tab=t.0#heading=h.gfwh8242orfp)  
- 📑 [API & Service Specification](https://editor.swagger.io/?url=https://raw.githubusercontent.com/egovernments/DIGIT-Studio/refs/heads/master/design/generic-service.yaml)  
- ⚙️ [Sample Service Configuration](./design/serviceConfig.json)

---
