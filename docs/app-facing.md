# Application-/Backend-facing API
The Application- or Backend-facing REST API is the one that is only accessible by your application's backend. It is not accessible directly from your frontend or the internet. The connection between the REST API Server and your backend which invoked the HTTP REST calls is authenticated and protected using mutual TLS (mTLS).

## Create user
Create a new user.

URL: ```/users/```

Method: ```POST```

JSON Payload: 
```
{
    "email": "<user's email address>",
    "password": "<User's password (min length = 8, max  length = 32)>",
    "confirmed": true|false,
    "enabled": true|false,
    "data": {}
}
```

HTTP Response Status Codes:

* 201: Created (user successfully created, User ID in response header 'X-Object-ID')
* 400: Bad request (invalid JSON payload)
* 409: Conflict (email address already exists)

## Get user
Get a user object.

URL: ```/users/<ID>```

Method: ```GET```

HTTP Response Status Codes:

* 200: OK (successful, result in response body payload)
* 404: Not found (invalid User ID)

HTTP Response Body:
```
{
    "email": "<user's email address>",
    "password": "<User's password (min length = 8, max  length = 32)>",
    "confirmed": true|false,
    "enabled": true|false,
    "data": {}
}
```

## Delete user
Delete a user.

URL: ```/users/<ID>```

Method: ```DELETE```

HTTP Response Status Codes:

* 204: No content (successful)
* 404: Not found (invalid User ID)

## Set email address
Set a user's email address.

URL: ```/users/<ID>/email```

Method: ```PUT```

JSON Payload: 
```
{
    "email": "<user's email address>"
}
```

HTTP Response Status Codes:

* 204: No content (successful)
* 400: Bad request (invalid JSON payload)
* 404: Not found (invalid User ID)
* 409: Conflict (email address already exists)

## Set password
Set a user's password.

URL: ```/users/<ID>/password```

Method: ```PUT```

JSON Payload: 
```
{
    "password": "<user's new password>"
}
```

HTTP Response Status Codes:

* 204: No content (successful)
* 400: Bad request (invalid JSON payload)
* 404: Not found (invalid User ID)

## Disable user
Disable a user account so that the user can't log in anymore.

URL: ```/users/<ID>/disable```

Method: ```PUT```

HTTP Response Status Codes:

* 204: No content (successful)
* 404: Not found (invalid User ID)

## Enable user
Enable a user account so that the user can log in.

URL: ```/users/<ID>/enable```

Method: ```PUT```

HTTP Response Status Codes:

* 204: No content (successful)
* 404: Not found (invalid User ID)

## Set custom user data
Store custom JSON data in a user object.

URL: ```/users/<ID>/data```

Method: ```PUT```

JSON Payload: 
```
{
    <Custom JSON data>
}
```

HTTP Response Status Codes:

* 204: No content (successful)
* 400: Bad request (invalid JSON payload)
* 404: Not found (invalid User ID)

## Get custom user data
Retrieve previously stored custom JSON data from a user object.

URL: ```/users/<ID>/data```

Method: ```GET```

HTTP Response Status Codes:

* 204: No content (successful)
* 404: Not found (invalid User ID)

HTTP Response Body:
```
{
    <Custom JSON data>
}
```

## Check password
Checks if a supplied plain-text password matched the user's hashed password.

URL: ```/users/<ID>/checkpw```

Method: ```POST```

JSON Payload: 
```
{
    "password": "<plain-text password>"
}
```

HTTP Response Status Codes:

* 200: OK (successful)
* 400: Bad request (invalid JSON payload)
* 404: Not found (invalid User ID)
  
HTTP Response Body:
```
{
    "result": true|false
}
```
