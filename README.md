# JWT Auth Proxy
This JWT Auth Proxy is a lightweight authentication proxy written in Go designed for use in Docker/Kubernetes environments.

![JWT Auth Proxy](https://raw.githubusercontent.com/virtualzone/jwt-auth-proxy/master/diagram.png)

JWT Auth Proxy sits between your frontend and your application's backend, handles authentication, and proxies authenticated requests to your backend. It offers REST APIs and logic for signup (incl. double-opt-in), password reset ("forgot my password") and verified email address change. This way, your application's backend can focus on the actual business logic, while relying on secure authentication having been performed before.

To your application's backend, the JWT Auth Proxy provides an mTLS-secured REST API for modifying user objects and storing custom data per user.

JWT Auth Proxy uses short-lived JWT access tokens (HMAC-signing with SHA-512) and long-lived UUIDv4 refresh tokens for securely retrieving new access tokens before the old one expires. It supports Two-Factor Authentication (2FA) via Time-based One-Time passwords (TOTP).

## Features
### User-facing
* Easy-to-use REST API for
  * Signup with double-opt-in
  * Login (with TOTP optionally)
  * Logout
  * Password reset (forgot password)
  * Email address change with double-opt-in
  * JWT access token renewal using long-lived refresh tokens
  * Activate and disable Two-Factor Authentication (2FA, TOTP)
* Proxy authenticated requests to your application's backend
* Whitelist backend URLs not requiring authentication (or blacklist)

### Application-/Backend-facing
* mTLS encrypted connecting (mutual TLS)
* Easy-to-use REST API for
  * Create user
  * Delete user
  * Disable/enable user
  * Check password
  * Set password
  * Set email address
  * Store and retrieve custom per-user data (JSON)

## Example
There is a little sample application in the [example folder](https://github.com/virtualzone/jwt-auth-proxy/tree/master/example). It consists of the JWT Auth Proxy, a React-based web frontend and a Go-based application backend. It also contains a MongoDB instance and an instance of the Mailhog fake SMTP server.

To start the example:

```
git clone https://github.com/virtualzone/jwt-auth-proxy.git
cd jwt-auth-proxy/example
docker-compose up -d
```

Access the frontend at: http://localhost:8080

To access the Mailhog frontend and check for your signup mail: http://localhost:8025

## Setup
Prerequisites: A MongoDB instance (tested with MongoDB v4).

Please refer to the [docker-compose.yml](https://github.com/virtualzone/jwt-auth-proxy/blob/master/example/docker-compose.yml) example on how to use the pre-build Docker image.

The server requires a key and a certificate for providing mTLS encryption for the backend REST API.

You can generate the server's and clients' keys and certificates manually using OpenSSL. By default, the server will create a CA and generate keys and certificates for both server and client on startup (see BACKEND_GENERATE_CERT environment variable).

## Handling proxied requests in your application's backend
When your application's backend receives an HTTP request proxied through the JWT Auth Proxy, it receives all the HTTP request headers sent by the HTTP client/browser, plus:

* ```Authorization```: The successfully validated JWT access token (format: ```Bearer <Token>```).
* ```X-Auth-UserID```: The user's ID you can use to make calls to the backend-facing REST API.
* ```Forwarded```: Information from the client-facing side of the proxy server.
* ```X-Forwarded-For``` (XFF): The originating IP address of the client.
* ```X-Forwarded-Host``` (XFH): The original host requested by the client in the Host HTTP request header.
* ```X-Forwarded-Proto``` (XFP): The protocol (HTTP or HTTPS) the client used to connect to the proxy.

## User-facing REST API
### Sign up / register new user
* Use Case: Sign up a new user using his unique email address as the username.
* URL: /auth/signup
* Method: POST
* JSON Payload: 
  ```
  {
      "email": "<User's email address = username>",
      "password": "<User's chosen password (min length = 8, max  length = 32)>"
  }
  ```
* HTTP Response Status Codes:
  * 201: Created (user successfully signed up, User ID in response header 'X-Object-ID')
  * 400: Bad request (invalid JSON payload)
  * 409: Conflict (user already exists)

### Log in
* Use case: Log in an activated and enabled user, retrieve Access and Refresh Tokens.
* URL: /auth/login
* Method: POST
* JSON Payload: 
  ```
  {
      "email": "<User's email address = username>",
      "password": "<User's chosen password (min length = 8, max  length = 32)>",
      "otp": "<Six digit TOTP>"
  }
  ```
* HTTP Response Status Codes:
  * 200: OK (user successfully logged in or additional TOTP required, result in response body payload)
  * 400: Bad request (invalid JSON payload)
  * 401: Unauthorized (authorization failed due to various reasons)
* HTTP Response Body for successful login:
  ```
  {
      "accessToken": "<short-lived JWT Access Token>",
      "refreshToken": "<long-lived UUIDv4 Refresh Token>",
  }
  ```
* HTTP Response Body if TOTP is required:
  ```
  {
      "otpRequired": true
  }
  ```

### Refresh Access Token
* Use case: Refresh short-lived Access Token with long-lived Refresh Token.
* URL: /auth/refresh
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* JSON Payload: 
  ```
  {
      "refreshToken": "<long-lived UUIDv4 Refresh Token from login>"
  }
  ```
* HTTP Response Status Codes:
  * 200: OK (successful, result in response body payload)
  * 400: Bad request (invalid JSON payload)
  * 401: Unauthorized (authorization failed due to various reasons)
* HTTP Response Body:
  ```
  {
      "accessToken": "<short-lived JWT Access Token>",
      "refreshToken": "<long-lived UUIDv4 Refresh Token>",
  }
  ```

### Log out
* Use case: Invalidate Refresh Token.
* URL: /auth/logout
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* JSON Payload: 
  ```
  {
      "refreshToken": "<long-lived UUIDv4 Refresh Token from login>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 400: Bad request (invalid JSON payload)
  * 401: Unauthorized (authorization failed due to various reasons)

### Ping
* Use case: Check if Access Token is still valid.
* URL: /auth/ping
* Method: GET
* Request Header: ```Authorization: Bearer <Access Token>```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 401: Unauthorized (authorization failed due to various reasons)

### Confirm
* Use case: User wants to confirm a requests received via email (such as signup, password reset, email change)
* URL: /auth/confirm/```<ID from email>```
* Method: POST
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 404: Not found (invalid, expired or already confirmed ID)

### Set password
* Use case: Logged in user wants to change his password.
* URL: /auth/setpw
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* JSON Payload: 
  ```
  {
      "oldPassword": "<user's old password>",
      "newPassword": "<user's new password>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 400: Bad request (invalid JSON payload)
  * 401: Unauthorized (authorization failed due to various reasons)

### Change email address
* Use case: Logged in user wants to change his email address (= username).
* URL: /auth/changeemail
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* JSON Payload: 
  ```
  {
      "password": "<user's password>",
      "email": "<user's new email address>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful, email sent to new email address - confirmation required before new address gets activated)
  * 400: Bad request (invalid JSON payload)
  * 401: Unauthorized (authorization failed due to various reasons)
  * 409: Conflict (email address already exists)

### Reset password
* Use case: User forgot his password and wants to reset it.
* URL: /auth/initpwreset
* Method: POST
* JSON Payload: 
  ```
  {
      "email": "<user's email address>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful, email sent user - confirmation required before new password is generated)
  * 400: Bad request (invalid JSON payload)

### Delete account
* Use case: User wants to delete his own account.
* URL: /auth/delete
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* JSON Payload: 
  ```
  {
      "password": "<user's password>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 400: Bad request (invalid JSON payload)
  * 401: Unauthorized (authorization failed due to various reasons)

### TOTP Initialization
* Use case: User wants activate Time-bases One-Time passwords (TOTP) for his account.
* URL: /auth/otp/init
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* HTTP Response Status Codes:
  * 200: OK (successful, result in response body payload)
  * 400: Bad request (i.e. TOTP already activated)
* HTTP Response Body:
  ```
  {
      "secret": "<TOTP secret>",
      "image": "<Base64 encoded PNG image of QR code>",
  }
  ```

### TOTP Confirmation
* Use case: User wants to confirm TOTP activation after scanning the previously generated QR Code with his authenticator app.
* URL: /auth/otp/confirm
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* JSON Payload: 
  ```
  {
      "passcode": "<Six digit TOTP>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 400: Bad request (i.e. invalid TOTP)

### TOTP Deactivation
* Use case: User wants disable TOTP.
* URL: /auth/otp/disable
* Method: POST
* Request Header: ```Authorization: Bearer <Access Token>```
* HTTP Response Status Codes:
  * 204: No content (successful)

## Application-/Backend-facing REST API
### Create user
* Use case: Create a new user.
* URL: /users/
* Method: POST
* JSON Payload: 
  ```
  {
      "email": "<user's email address>",
      "password": "<User's password (min length = 8, max  length = 32)>",
      "confirmed": true|false,
      "enabled": true|false,
      "data": {}
  }
  ```
* HTTP Response Status Codes:
  * 201: Created (user successfully created, User ID in response header 'X-Object-ID')
  * 400: Bad request (invalid JSON payload)
  * 409: Conflict (email address already exists)

### Get user
* Use case: Get a user object.
* URL: /users/```<ID>```
* Method: GET
* HTTP Response Status Codes:
  * 200: OK (successful, result in response body payload)
  * 404: Not found (invalid User ID)
* HTTP Response Body:
  ```
  {
      "email": "<user's email address>",
      "password": "<User's password (min length = 8, max  length = 32)>",
      "confirmed": true|false,
      "enabled": true|false,
      "data": {}
  }
  ```

### Delete user
* Use case: Delete a user.
* URL: /users/```<ID>```
* Method: DELETE
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 404: Not found (invalid User ID)

### Set email address
* Use case: Set a user's email address.
* URL: /users/```<ID>```/email
* Method: PUT
* JSON Payload: 
  ```
  {
      "email": "<user's email address>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 400: Bad request (invalid JSON payload)
  * 404: Not found (invalid User ID)
  * 409: Conflict (email address already exists)

### Set password
* Use case: Set a user's password.
* URL: /users/```<ID>```/password
* Method: PUT
* JSON Payload: 
  ```
  {
      "password": "<user's new password>"
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 400: Bad request (invalid JSON payload)
  * 404: Not found (invalid User ID)

### Disable user
* Use case: Disable a user account so that the user can't log in anymore.
* URL: /users/```<ID>```/disable
* Method: PUT
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 404: Not found (invalid User ID)

### Enable user
* Use case: Enable a user account so that the user can log in.
* URL: /users/```<ID>```/enable
* Method: PUT
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 404: Not found (invalid User ID)

### Set custom user data
* Use case: Store custom JSON data in a user object.
* URL: /users/```<ID>```/data
* Method: PUT
* JSON Payload: 
  ```
  {
      <Custom JSON data>
  }
  ```
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 400: Bad request (invalid JSON payload)
  * 404: Not found (invalid User ID)

### Get custom user data
* Use case: Retrieve previously stored custom JSON data from a user object.
* URL: /users/```<ID>```/data
* Method: GET
* HTTP Response Status Codes:
  * 204: No content (successful)
  * 404: Not found (invalid User ID)
* HTTP Response Body:
  ```
  {
      <Custom JSON data>
  }
  ```

### Check password
* Use case: Checks if a supplied plain-text password matched the user's hashed password.
* URL: /users/```<ID>```/checkpw
* Method: POST
* JSON Payload: 
  ```
  {
      "password": "<plain-text password>"
  }
  ```
* HTTP Response Status Codes:
  * 200: OK (successful)
  * 400: Bad request (invalid JSON payload)
  * 404: Not found (invalid User ID)
* HTTP Response Body:
  ```
  {
      "result": true|false
  }
  ```

## Configuration Options
Configuration is performed by setting the appropriate environment variables.

Env | Default | Description
--- | --- | ---
JWT_SIGNING_KEY | 32 Bytes Random String | The private key for signing the JWT access tokens.
PUBLIC_LISTEN_ADDR | 0.0.0.0:8080 | The listening address for the user-facing HTTP server.
PUBLIC_API_PATH | /auth/ | The path for the user-facing REST API.
BACKEND_LISTEN_ADDR | 0.0.0.0:8443 | The listening address for the backend-facing HTTPS server.
BACKEND_CERT_DIR | ./certs/ | The directory containing the backend-facing HTTP server's certificates (mTLS).
BACKEND_GENERATE_CERT | 1 | Whether to create CA and server key-pair on startup (= 1).
BACKEND_CERT_HOSTNAMES | localhost | The hostnames to generate the server certificate for, separated by commas.
BACKEND_CERT_IPS | 127.0.0.1,::1 | The IP addresses to generate the server certificate for, separated by commas.
TEMPLATE_SIGNUP | res/signup.tpl | The email template for signup confirmation mails.
TEMPLATE_CHANGE_EMAIL | res/changeemail.tpl | The email template for email address change confirmation mails.
TEMPLATE_RESET_PASSWORD | res/resetpassword.tpl | The email template for password reset confirmation mails.
TEMPLATE_NEW_PASSWORD | res/newpassword.tpl | The email template for new password mails.
MONGO_DB_URL | mongodb://localhost:27017 | The URL of the MongoDB database server.
MONGO_DB_NAME | jwt_auth_proxy | The database name of the MongoDB database.
CORS_ENABLE | 0 | Whether to enable (= 1) Cross-Origin Resource Sharing (CORS) response headers.
CORS_ORIGIN | * | The value of the 'Access-Control-Allow-Origin' header.
CORS_HEADERS | * | The value of the 'Access-Control-Allow-Headers' header.
SMTP_SERVER | 127.0.0.1:25 | The address and port of the outgoing SMTP server.
SMTP_SENDER_ADDR | no-reply@localhost | The SMTP sender address.
ALLOW_SIGNUP | 1 | Whether to allow (= 1) signup requests at the user-facing HTTP server.
ALLOW_CHANGE_PASSWORD | 1 | Whether to allow (= 1) change password requests at the user-facing HTTP server.
ALLOW_CHANGE_EMAIL | 1 | Whether to allow (= 1) change email address requests at the user-facing HTTP server.
ALLOW_FORGOT_PASSWORD | 1 | Whether to allow (= 1) password reset requests at the user-facing HTTP server.
ALLOW_DELETE_ACCOUNT | 1 | Whether to allow (= 1) "delete my account" requests at the user-facing HTTP server.
TOTP_ENABLE | 0 | Whether to enable (= 1) support for Time-based One-Time Passwords (TOTP) as a second authentication factor (2FA).
TOTP_ISSUER | JWT Auth Proxy | The TOTP Issuer.
TOTP_ENCRYPT_KEY | '' | The passphrase encrypt the TOTP Secrets in the database (minimum length: 16 bytes). Required if TOTP_ENABLE=1.
PROXY_TARGET | http://127.0.0.1:80 | The target server hosting your application backend.
PROXY_WHITELIST | '' | Whitelisted URL prefixes at the target server not requiring a valid authentication. Separate prefixes by colons (':'). Don't use with PROXY_BLACKLIST.
PROXY_BLACKLIST | '' | Blacklisted URL prefixes at the target server requiring a valid authentication. Separate prefixes by colons (':'). Don't use with PROXY_WHITELIST.
ACCESS_TOKEN_LIFETIME | 5 | The access token lifetime in minutes.
REFRESH_TOKEN_LIFETIME | 1,440 | The refresh token lifetime in minutes.
PENDING_ACTION_LIFETIME | 1,440 | The lifetime of pending actions (such as confirmation requests) in minutes.
