# IAM Test Framework

## Overview
The IAM Test Framework is a generic testing framework designed to validate the basic functionalities of various Identity and Access Management (IAM) solutions, including Zitadel, Keycloak, and Casdoor. It focuses on standard API interactions for user management, authentication, and authorization flows.

## Features
- User Registration
- User Login
- OIDC Authorization Flow
- Token Issuance and Validation
- User Info Retrieval
- Logout Functionality

## Requirements
- Go 1.16 or higher
- Docker (for containerization)
- Kubernetes (for deployment)

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/iam-test-framework.git
cd iam-test-framework