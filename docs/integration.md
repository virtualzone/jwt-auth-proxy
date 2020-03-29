# Application Integration
JWT Auth Proxy usually integrates with your application like this:

* Your frontend (i.e. JavaScript based) invokes the User-facing REST API to perform signups, validate logins requests, invoke password resets and more.
* Your application's backend (i.e. REST services) is "hidden" behind the JWT Auth Proxy, receiving only correctly authenticated requests or invokations to URLs on the proxy's whitelist.

## HTTP Request Headers
When your application's backend receives an HTTP request proxied through the JWT Auth Proxy, it receives all the HTTP request headers sent by the HTTP client/browser, plus:

* ```Authorization```: The successfully validated JWT access token (format: ```Bearer <Token>```).
* ```X-Auth-UserID```: The user's ID you can use to make calls to the backend-facing REST API.
* ```Forwarded```: Information from the client-facing side of the proxy server.
* ```X-Forwarded-For``` (XFF): The originating IP address of the client.
* ```X-Forwarded-Host``` (XFH): The original host requested by the client in the Host HTTP request header.
* ```X-Forwarded-Proto``` (XFP): The protocol (HTTP or HTTPS) the client used to connect to the proxy.

## Calling the Backend API
To call the backend-facing API, invoke REST-based HTTP requests from your backend to JWT Auth Proxy's backend-facing REST service. This service is usually listening on port 8443 and requires a valid mTLS certificate.