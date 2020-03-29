# Setup
Prerequisites: A MongoDB instance (tested with MongoDB v4).

Please refer to the [docker-compose.yml](https://github.com/virtualzone/jwt-auth-proxy/blob/master/example/docker-compose.yml) example on how to use the pre-build Docker image.

The server requires a key and a certificate for providing mTLS encryption for the backend REST API.

You can generate the server's and clients' keys and certificates manually using OpenSSL. By default, the server will create a CA and generate keys and certificates for both server and client on startup (see BACKEND_GENERATE_CERT environment variable).