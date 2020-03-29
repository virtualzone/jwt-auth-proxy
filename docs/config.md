# Configuration
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