# Welcome to JWT Auth Proxy

## About JWT Auth Proxy
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