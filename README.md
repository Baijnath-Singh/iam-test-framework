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
  "iam_solutions": {
    "zitadel": {
      "domain": "http://127.0.0.1.sslip.io:8080",
      "registration_endpoint": "/v2/users/human",
      "login_endpoint": "/oauth/v2/login",
      "authorize_endpoint": "/oauth/v2/authorize",
      "token_endpoint": "/oauth/v2/token",
      "userinfo_endpoint": "/v2/users",
      "logout_endpoint": "/oauth/v2/logout",
      "token_request_params": {
        "client_credentials": {
          "grant_type": "client_credentials",
          "client_id": "nvo",
          "client_secret": "xVERaGld3aBl87plN5Zo1QdwCAaADyOYHLhX8zT5ip3WDZwtYJ28GGXXXXXXXXXX",
          "scope": "openid profile urn:zitadel:iam:org:project:id:zitadel:aud"
        }
      }
    },
    "keycloak": {
      "domain": "http://your-keycloak-domain.com",
      "registration_endpoint": "/auth/realms/{realm}/protocol/openid-connect/registration",
      "login_endpoint": "/auth/realms/{realm}/protocol/openid-connect/token",
      "authorize_endpoint": "/auth/realms/{realm}/protocol/openid-connect/auth",
      "token_endpoint": "/auth/realms/{realm}/protocol/openid-connect/token",
      "userinfo_endpoint": "/auth/realms/{realm}/protocol/openid-connect/userinfo",
      "logout_endpoint": "/auth/realms/{realm}/protocol/openid-connect/logout",
      "token_request_params": {
        "client_credentials": {
          "grant_type": "client_credentials",
          "client_id": "your-keycloak-client-id",
          "client_secret": "your-keycloak-client-secret",
          "scope": "openid"
        }
      }
    },
    "casdoor": {
      "domain": "http://your-casdoor-domain.com",
      "registration_endpoint": "/api/door/register",
      "login_endpoint": "/api/door/login",
      "authorize_endpoint": "/api/door/authorize",
      "token_endpoint": "/api/door/token",
      "userinfo_endpoint": "/api/door/userinfo",
      "logout_endpoint": "/api/door/logout",
      "token_request_params": {
        "client_credentials": {
          "grant_type": "client_credentials",
          "client_id": "your-casdoor-client-id",
          "client_secret": "your-casdoor-client-secret",
          "scope": "openid"
        }
      }
    }
  },
  "default_solution": "zitadel"
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

## How to test SSO/SLO
Register web application for App1
Register web application for App2
Clone https://github.com/Baijnath-Singh/iam-test-framework/tree/master/zitadel-sso-slo-test
In two separate terminals run App1 and App2 <go run main.go>
Open url http://localhost:3001 and http://localhost:3002 for App1 and App2 respectively
Make sure we are running the App first time or clear the cookies
For SSO
Click on login  link of App1, you would be redirected to login page where cred needs to be entered. 
After successful login, user information along with logout link is displayed on browser.
Click login link of App2  You would not redirected to login page as there is session active with App1 login. 
User information along with logout link is displayed on browser.
For SLO
Click logout link either of App1 or App2
After clicking logout link of say App1, you would be redirected to logout page
Try login to App2 now, you would be redirected to singed out page and will have to login to App2
