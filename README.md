# 🚀 DIGIT Studio ( Service Delivery Framework )

The **DIGIT Service Delivery Framework** is a **low-code/no-code** platform built to help government agencies and partners **rapidly design, configure, and deploy** digital public services—such as trade licenses, permits, and grievance redressals—with **minimal engineering effort**.

It builds on top of the proven DIGIT Core Platform and enables the rollout of fully functional digital services through **configuration** (no code) and **extensions** (code where needed).

---

## 🔍 Overview

### ✅ What Problem Does It Solve?

Traditionally, each government service is developed from scratch, which leads to:

- 🛠️ Redundant development work  
- 🔄 Inconsistent logic and user experience  
- 🕒 Longer timelines and higher implementation costs  

The **DIGIT Service Delivery Framework** solves this by providing a **shared foundation** where services can be **built once and reused across departments**, improving speed, cost, and quality of delivery.

---

## 🧱 Framework Components


### 1️⃣ Service Designer

A **visual low-code/no-code tool** for administrators and business users to define digital service configurations:

- Dynamic forms and field validations  
- Workflow steps and role permissions  
- Required documents and fee logic  

📝 The configurations are exported as **JSON** and stored in the Service Registry.

---

### 2️⃣ Service Registry

The **source of truth** for all configured services. It:

- Stores metadata, forms, workflows, roles, etc.  
- Serves configurations to the Service Runtime at runtime  

This allows services to be updated or versioned without code changes.

---

### 3️⃣ Service Runtime

The **execution engine** that:

- Renders forms dynamically  
- Manages workflows and state transitions  
- Validates inputs and uploads documents  
- Orchestrates API calls to DIGIT Core Services  

It consumes configurations from the Service Registry and brings them to life.

---

### 4️⃣ Service UI

A **dynamic front-end layer** that adapts based on the definitions from the Service Designer. It allows:

- 📥 Citizens and businesses to apply for services  
- 🧑‍💼 Employees to review and approve requests  
- 🖥️ A unified and responsive experience across devices  

---

### 5️⃣ Service Extensions

Plug-in components developed by developers to:

- Implement complex business logic  
- Extend default behavior (e.g., custom validations or fee rules)  
- Integrate with external systems  

These are useful where configuration alone is not enough.

---

### 6️⃣ Service Analytics

Built on **ElasticSearch**, it:

- Captures system events, workflow states, and user actions  
- Supports performance monitoring and reporting  
- Enables deep insights into service usage

---

### 7️⃣ Service Dashboard

Powered by **Kibana**, it offers:

- 📊 Visual dashboards for administrators  
- 📈 Real-time monitoring and usage analytics  
- 🧩 Configurable views by service or user role

---

## 🔗 Integration with DIGIT Core Services

The framework interacts with DIGIT's reusable backend services, including:

- 🔁 Workflow Management  
- 🔐 Authentication & Authorization  
- 🌐 Localization & MDMS (Master Data Management)  
- 📍 Location & Notification Services  
- 🗃️ Document Management, Audit Logging, and Persistence  

These APIs ensure **consistency**, **interoperability**, and **security** across all services.

---

## 👥 User Roles

Different users interact with the system through a unified DIGIT UI, personalized to their roles:

| Role               | Responsibilities                                   |
|--------------------|----------------------------------------------------|
| Service Designers  | Define service configurations                      |
| Citizens & Businesses | Consume services, submit applications             |
| Employees & Vendors | Review, verify, and approve applications          |
| Administrators     | Monitor operations, performance, and service usage |

---

## 🎯 Benefits

- ⚡ **Fast Deployment** – Roll out services in **days**, not months  
- 🧩 **Consistent Experience** – Standardized UI & logic across services  
- 🛠️ **Flexible Customization** – Use configuration or code where needed  
- 📈 **Scalable Architecture** – Add new services without core changes  
- 👨‍👩‍👧‍👦 **Citizen-Centric** – Unified access, tracking, and engagement  

---

## 🏗️ Ideal For

- Urban and rural service delivery  
- Permits, licenses, and benefit schemes  
- Grievance redressal and citizen engagement  
- Internal government workflows and automation

---

## 🧱 Repository Structure
  ```
  digit-studio/
  ├── design/                         # System design and configuration specs
  │   ├── design.md                   # Architecture & design overview
  │   ├── serviceConfig.json          # Example of a configurable service
  │   └── generic-service.yaml        # Service specifications
  │
  ├── frontend/                       # Frontend-related modules
  │   ├── common-ui/                  # Shared UI components and utilities
  │   └── console-ui/                 # Admin & control center UI
  │
  ├── backend/                        # Backend service modules
  │   ├── generic-service/            # Common backend utilities and logic
  │   ├── public-service/             # Backend powering dynamic UI components
  │   └── transformer/                # Admin & control processing layer
  │
  └── README.md                       # Project overview and documentation
  
  ```

---

## 📌 Documentation

  - 📐 [Design Document](https://docs.google.com/document/d/13LR7TQMsIg0nD5-Wdl4kj1r3kYjzLyKD0FVzvJkkR3s/edit?tab=t.0#heading=h.gfwh8242orfp)  

  - 📑 [API & Service Specification](https://editor.swagger.io/?url=https://raw.githubusercontent.com/egovernments/DIGIT-Studio/refs/heads/master/design/generic-service.yaml)  

  - ⚙️ [Sample Service Configuration](./design/serviceConfig.json)

---

## 🚀 Getting Started
  Clone the repository:

  ```bash
  git clone https://github.com/egovernments/DIGIT-Studio.git
  ```
  
  ```bash 
  cd DIGIT-Studio
  ```
Use docs/service-config.yaml as a reference to plug in your own services or UI variations.

---

## 🧩 Use Cases
  Unified master data management
  
  Configurable form-based workflows
  
  Service-level customization without redeployment
  
  Scalable support for new departments and use cases

---

## 🛠️ Version 1 Capabilities
  The first version of Digit Studio provides end-to-end capabilities for core service delivery and workflow-based applications, including:

Current Version Features:

| **Feature**                           | **Current Status**                                        | **Version** |
| ------------------------------------- | --------------------------------------------------------- | ----------- |
| **Apply**                             | 🟡 *In Progress – Positive flow implemented and deployed* |    v1       |
| **View & Workflow Transition**        | 🟡 *In Progress – Positive flow implemented and deployed* |    v1       |
| **Inbox**                             | 🔄 *In Progress – Pending deployment*                     |    v1       |
| **Search**                            | 🔄 *In Progress – Pending deployment*                     |    v1       |
| **Other Misc. (e.g., Tenant Config)** | 🔄 *In Progress – Pending deployment*                     |    v1       |
| **Applicant – Individual Support**    | 🚫 *Not started*                                          |    v1       |
| **Applicant – Organization Support**  | 🚫 *Not started*                                          |    v2       |
| **Bill & Payment**                    | 🚫 *Not started*                                          |    v1       |
| **PDF Generation & Integration**      | 🚫 *Not started*                                          |    v1       |
| **SMS Integration**                   | 🚫 *Not started*                                          |    v1       |
| **Checklist Integration**             | 🚫 *Not started*                                          |    v1       |
| **Edit & Resubmit Application**       | 🚫 *Not started*                                          |    v2       |
| **User Type Enablement**              | 🚫 *Not started*                                          |    v2       |
| **Service Initialization**            | 🚫 *Not started*                                          |    v2       |
| **Console**                           | 🚫 *Not started*                                          |    v3       |


### Status Legend:
  
✅ **Completed**                           
🟡 **In Progress – Positive Flow Handled & Deployed**
🔄 **In Progress – Not Deployed**  
🔒 **In Progress – Blocked**   
🚫 **Not Started**

---

## 🧩 Key Features

### ✅ Configuration via Service Designer *(Planned via UI, Manual in Alpha)*

Administrators define and manage the service using configurations for:

- 📄 Application forms and field validations  
- 📎 Required documents  
- 💸 Fee calculation rules  
- 🧑‍💼 Role-based access control and workflow steps  
- 🔔 Notification triggers (SMS, Email)  

> ℹ️ In the alpha release, these configurations are authored manually in JSON format.

---

### 👥 Citizen Interaction through Dynamic UI

Citizens or business users access the service through a **dynamic UI**, rendered based on configuration:

- 📝 Fill and submit application forms  
- 📎 Upload required documents  
- 💳 Make secure online payments  
- 🌐 Multilingual and mobile responsive interface  

---

### 🏛️ Employee Processing via Workflow

Municipal staff access the same platform to:

- 🔍 Review and verify applications  
- 🗂️ Approve or reject requests  
- 🧭 Follow configured role-based workflows  
- ❓ Raise clarifications with pre-defined reasons  

Workflow orchestration is managed via the **Workflow Core Service** with automatic task assignments.

---

### ⚙️ Backend Orchestration via Service Runtime

The **Service Runtime** handles:

- ✅ Input validation (schema + custom validators)  
- 🔁 Workflow transitions and audit logging  
- 💸 Billing and payment integration  
- 📎 File uploads and notifications  
- 📜 PDF certificate generation  

All backend services are composed through DIGIT’s modular and reusable service components.

---

### 📊 Monitoring & Analytics

- Metrics such as application status, turnaround times, and usage patterns are captured via **Service Analytics**.  
- Visualization is supported through the **Kibana Dashboard** for monitoring and operational insights.

---

## 🚀 Alpha Release Highlights

The **Alpha version** delivers a functional runtime and service configuration model with the following focus:

### 🎯 Included Components

- **Service Runtime**  
  - Executes services using JSON configurations  
  - Handles dynamic form rendering, workflow, billing, and notifications  
- **Core Integrations**  
  - Workflow, Billing, FileStore, Notification, Localization, MDMS  
  - ElasticSearch-based Service Analytics  

### 🧾 Manual Configuration Setup

- Services are configured via JSON files following a structured schema  
- Includes: forms, validations, workflow steps, roles, fees, etc.  
- Versioned **Service Registry** is used to store and serve configurations  

---

## ⚠️ Known Limitations

| Area | Limitation |
|------|------------|
| **Standard UI Flow** | Only supports Apply → Workflow → Approvals |
| **Validation** | Only basic schema validations supported; complex logic requires custom APIs |
| **Billing Logic** | Must be implemented separately; not configurable |
| **UI Schema** | No field-level dynamic behaviors or custom UI schemas |
| **Form Structure** | Only supports one level of nesting; complex tables not supported |
| **Search** | JSON-based storage limits advanced querying |
| **Workflow** | No support for parallel workflows |
| **Service Designer** | Visual UI not available in Alpha |
| **Service Initializer** | May require manual prefill; automation to be added later |
| **Import/Export** | Not supported for service configs or data |
| **Registry Integration** | External registry sync not supported |
| **Applicant Type** | Only Individual applicants supported; organizations not supported yet |

---



## 🔭 Roadmap: Future Enhancements
  In future versions, we aim to address:
  
  🧑‍💻 Enhanced Admin Console for better role-based access and control
  
  🏢 External Registries Integration for real-time data exchange
  
  🔄 Additional use cases across departments with plug-and-play capabilities
  
  📊 Analytics & Reporting modules
  
  🧠 AI-assisted data suggestions and automation

---


## 🤝 Contributing
  Contributions are welcome! Please refer to the contributing guide for guidelines on submitting issues or pull requests.

## 📬 Contact
  For any questions or support, reach out to [jagankumar](https://github.com/jagankumar-egov)

## 🛡️ License
  This project is licensed under the MIT License. See the LICENSE file for details.

