# User-facing API
The User- or Frontend-facing REST API is the one that is directly accessible by your users. The REST API is usually invoked by JavaScript code that ships with your web frontend.

## Sign up / register new user
Sign up a new user using his unique email address as the username.

URL: ```/auth/signup```

Method: ```POST```

JSON Payload: 
```
{
    "email": "<User's email address = username>",
    "password": "<User's chosen password (min length = 8, max  length = 32)>"
}
```
    
HTTP Response Status Codes:

* 201: Created (user successfully signed up, User ID in response header 'X-Object-ID')
* 400: Bad request (invalid JSON payload)
* 409: Conflict (user already exists)

## Log in
Log in an activated and enabled user, retrieve Access and Refresh Tokens.

URL: ```/auth/login```

Method: ```POST```

JSON Payload: 
```
{
    "email": "<User's email address = username>",
    "password": "<User's chosen password (min length = 8, max  length = 32)>",
    "otp": "<Six digit TOTP>"
}
```
HTTP Response Status Codes:

* 200: OK (user successfully logged in or additional TOTP required, result in response body payload)
* 400: Bad request (invalid JSON payload)
* 401: Unauthorized (authorization failed due to various reasons)

HTTP Response Body for successful login:
```
{
    "accessToken": "<short-lived JWT Access Token>",
    "refreshToken": "<long-lived UUIDv4 Refresh Token>",
}
```

HTTP Response Body if TOTP is required:
```
{
    "otpRequired": true
}
```

## Refresh Access Token
Refresh short-lived Access Token with long-lived Refresh Token.

URL: ```/auth/refresh```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

JSON Payload: 
```
{
    "refreshToken": "<long-lived UUIDv4 Refresh Token from login>"
}
```

HTTP Response Status Codes:
* 200: OK (successful, result in response body payload)
* 400: Bad request (invalid JSON payload)
* 401: Unauthorized (authorization failed due to various reasons)

HTTP Response Body:
```
{
    "accessToken": "<short-lived JWT Access Token>",
    "refreshToken": "<long-lived UUIDv4 Refresh Token>",
}
```

## Log out
Invalidate Refresh Token.

URL: ```/auth/logout```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

JSON Payload: 
```
{
    "refreshToken": "<long-lived UUIDv4 Refresh Token from login>"
}
```

HTTP Response Status Codes:
* 204: No content (successful)
* 400: Bad request (invalid JSON payload)
* 401: Unauthorized (authorization failed due to various reasons)

## Ping
Check if Access Token is still valid.

URL: ```/auth/ping```

Method: ```GET```

Request Header: ```Authorization: Bearer <Access Token>```

HTTP Response Status Codes:

* 204: No content (successful)
* 401: Unauthorized (authorization failed due to various reasons)

## Confirm
User wants to confirm a requests received via email (such as signup, password reset, email change)

URL: ```/auth/confirm/<ID from email>```

Method: ```POST```

HTTP Response Status Codes:

* 204: No content (successful)
* 404: Not found (invalid, expired or already confirmed ID)

## Set password
Logged in user wants to change his password.

URL: ```/auth/setpw```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

JSON Payload: 
```
{
    "oldPassword": "<user's old password>",
    "newPassword": "<user's new password>"
}
```

HTTP Response Status Codes:

* 204: No content (successful)
* 400: Bad request (invalid JSON payload)
* 401: Unauthorized (authorization failed due to various reasons)

## Change email address
Logged in user wants to change his email address (= username).

URL: ```/auth/changeemail```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

JSON Payload: 
```
{
    "password": "<user's password>",
    "email": "<user's new email address>"
}
```

HTTP Response Status Codes:

* 204: No content (successful, email sent to new email address - confirmation required before new address gets activated)
* 400: Bad request (invalid JSON payload)
* 401: Unauthorized (authorization failed due to various reasons)
* 409: Conflict (email address already exists)

## Reset password
User forgot his password and wants to reset it.

URL: ```/auth/initpwreset```

Method: ```POST```

JSON Payload: 
```
{
    "email": "<user's email address>"
}
```

HTTP Response Status Codes:

* 204: No content (successful, email sent user - confirmation required before new password is generated)
* 400: Bad request (invalid JSON payload)

## Delete account
User wants to delete his own account.

URL: ```/auth/delete```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

JSON Payload: 
```
{
    "password": "<user's password>"
}
```

HTTP Response Status Codes:

* 204: No content (successful)
* 400: Bad request (invalid JSON payload)
* 401: Unauthorized (authorization failed due to various reasons)

## TOTP Initialization
User wants activate Time-bases One-Time passwords (TOTP) for his account.

URL: ```/auth/otp/init```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

HTTP Response Status Codes:

* 200: OK (successful, result in response body payload)
* 400: Bad request (i.e. TOTP already activated)

HTTP Response Body:
```
{
    "secret": "<TOTP secret>",
    "image": "<Base64 encoded PNG image of QR code>",
}
```

## TOTP Confirmation
User wants to confirm TOTP activation after scanning the previously generated QR Code with his authenticator app.

URL: ```/auth/otp/confirm```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

JSON Payload: 
```
{
    "passcode": "<Six digit TOTP>"
}
```

HTTP Response Status Codes:

* 204: No content (successful)
* 400: Bad request (i.e. invalid TOTP)

## TOTP Deactivation
User wants disable TOTP.

URL: ```/auth/otp/disable```

Method: ```POST```

Request Header: ```Authorization: Bearer <Access Token>```

HTTP Response Status Codes:

* 204: No content (successful)