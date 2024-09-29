# IAM Test Framework

## Description

The IAM Test Framework is a Go-based application designed to test and validate various Identity and Access Management (IAM) solutions. It provides a streamlined interface for users to perform essential tasks such as token issuance and user information retrieval across different IAM providers, including Zitadel, Keycloak, and Casdoor.

## Framework Structure

The framework is organized into several directories, each serving a specific purpose:
```
iam-test-framework/
├── cmd/
│   └── main.go
├── api/
│   ├── iam_client.go
│   ├── requests.go
│   ├── responses.go
│   ├── error_handling.go
│   └── routes/
│       ├── registration.go
│       ├── login.go
│       ├── oidc.go
│       ├── token.go
│       ├── userinfo.go
│       └── logout.go
├── input/
│   └── user_input.go
├── config/
│   ├── config.go
│   └── config.json
├── logger/
│   └── logger.go
├── kubernetes/
│   └── deployment.yaml
└── README.md
```

## File Purpose

- **cmd/main.go**: The entry point of the application where user input is handled and IAM functionalities are invoked based on user selection.

- **api/iam_client.go**: Contains functions for making API calls to the IAM solutions.

- **api/requests.go**: Defines the request structures used for API interactions.

- **api/responses.go**: Defines the response structures received from IAM solutions.

- **api/error_handling.go**: Implements error handling mechanisms for API requests and responses.

- **api/routes/**: Contains individual route files for handling different IAM functionalities:
  - **registration.go**: Manages user registration processes.
  - **login.go**: Handles user login processes.
  - **oidc.go**: Manages OpenID Connect authorization flows.
  - **token.go**: Implements token issuance and management.
  - **userinfo.go**: Retrieves user information based on access tokens.
  - **logout.go**: Handles user logout processes.

- **input/user_input.go**: Contains functions to handle user input for selecting options and entering parameters.

- **config/config.go**: Loads and manages the configuration settings for the IAM solutions.

- **config/config.json**: JSON file containing configuration details for different IAM providers.

- **logger/logger.go**: Implements logging functionality for the application.

- **kubernetes/deployment.yaml**: Configuration file for deploying the application to Kubernetes.

## Config File Detail

The configuration for the IAM Test Framework is defined in `config/config.json`. This file includes:

- **IAMSolutions**: A JSON object where each key is the name of the IAM solution (e.g., Zitadel, Keycloak, Casdoor) and contains specific configuration details for that provider:
  - **Domain**: The base URL of the IAM provider.
  - **TokenEndpoint**: The endpoint for retrieving access tokens.
  - **UserinfoEndpoint**: The endpoint for fetching user information.
  - **TokenRequestParams**: Parameters required for token requests, including client ID, client secret, and scope.

### Example `config.json` Structure

```json
{
  "IAMSolutions": {
    "zitadel": {
      "Domain": "https://example.zitadel.ch",
      "TokenEndpoint": "/oauth/token",
      "UserinfoEndpoint": "/oauth/userinfo",
      "TokenRequestParams": {
        "client_id": "your_client_id",
        "client_secret": "your_client_secret",
        "scope": "openid profile"
      }
    },
  }
}
```

## How to Build
To build the IAM Test Framework, follow these steps:
1. Clone the repository:
git clone https://github.com/yourusername/iam-test-framework.git ; cd iam-test-framework

2. Install dependencies (if any):
go mod tidy

3. Build the application:
go build -o iam-test-framework cmd/main.go

## How to Use
1. Run the application:
./iam-test-framework

2. Follow the prompts to select the IAM functionality you wish to test:
User Registration
OIDC Authorization
Token Issuance
User Info Retrieval

3. Input the necessary parameters as prompted, based on the selected functionality.



