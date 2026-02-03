# CAS Authentication (for UAB)

This document outlines the design for adding CAS-based authentication to BlazeBoard, specifically for UAB's environment.

## 1. Protocol Flow

The authentication flow will follow the CAS 2.0 protocol:

1.  **Initiation:** A user clicks a "Sign in with UAB" button, which directs them to our `/cas/login` endpoint.
2.  **Redirect to CAS:** Our server redirects the user to the UAB CAS login page (`CAS_SERVER_URL`). A `service` parameter is included, pointing back to our callback URL (`/cas/callback`).
3.  **UAB Authentication:** The user authenticates with UAB's system (BlazerID, password, and any 2FA like Duo). This part is handled entirely by UAB's infrastructure.
4.  **Redirect to Callback:** Upon successful authentication, the CAS server redirects the user back to our `/cas/callback` endpoint with a short-lived `ticket`.
5.  **Ticket Validation:** Our backend makes a server-to-server GET request to the CAS `/serviceValidate` endpoint to validate the ticket.
6.  **User Provisioning:** The CAS validation response contains the user's username (BlazerID). We use this to provision an account:
    - **ID/Username:** The BlazerID itself.
    - **Email:** Derived as `<BlazerID>@uab.edu`.
    - **Name:** Defaults to the BlazerID (can be changed by the user later).
7.  **Session Creation:** A session cookie is created for the user, and they are redirected to the application's home page or their original destination.

## 2. Backend Components

### Environment Variables

A new configuration struct will be added to `app/pkg/env/env.go`.

```go
// app/pkg/env/env.go
type CAS struct {
    ServerURL  string `env:"CAS_SERVER_URL"`  // e.g., https://cas.uab.edu/cas
    ServiceURL string `env:"CAS_SERVICE_URL"` // Optional, defaults to app base URL
}
```

### CAS Package (`app/pkg/cas/cas.go`)

A new package to encapsulate CAS protocol logic.

-   `IsConfigured() bool`: Checks if `CAS_SERVER_URL` is set.
-   `LoginURL(redirectURL string) (string, error)`: Constructs the redirect URL for the CAS server.
-   `ValidateTicket(ticket string) (*Profile, error)`: Validates the ticket with the CAS server and parses the XML response.
-   `Profile`: A struct to hold the authenticated user's information (`ID`, `Email`, `Name`).

### Handlers (`app/handlers/cas.go`)

New HTTP handlers to manage the flow.

-   `CASLogin(c *fider.Context)`: Handles the initial login request and redirects to the CAS server.
-   `CASCallback(c *fider.Context)`: Handles the callback from the CAS server, validates the ticket, gets or creates the user account, and establishes a session. It will reuse the existing `actions.GetOrCreateUserFromProvider` action.

### Routes (`app/cmd/routes.go`)

New routes will be added to handle the CAS flow.

```go
// app/cmd/routes.go (inside the main route group)
r.Get("/cas/login", handlers.CASLogin())
r.Get("/cas/callback", handlers.CASCallback())
```

## 3. Frontend Components

### Sign-In Page (`public/pages/SignIn/SignIn.page.tsx`)

The sign-in page will be updated to conditionally display a "Sign in with UAB" button.

-   The backend will pass a `casEnabled: true` prop to the page if `CAS_SERVER_URL` is configured.
-   The button will link to `/cas/login?redirect=/`.
-   It will be styled with UAB's brand color for easy recognition.

## 4. Error Handling

-   **Invalid Ticket:** Redirect to the sign-in page with an error message.
-   **CAS Server Unreachable:** Log the error and show a generic failure message.
-   **User Not Invited (Private Tenant):** Redirect to a "not-invited" page, consistent with other providers.

## 5. User Identity

-   The provider name for these users will be `uab`, the same as the existing SAML implementation. This ensures that existing UAB users via SAML can sign in with CAS without creating a duplicate account.
