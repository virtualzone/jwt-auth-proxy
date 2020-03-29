# Setup

## Prerequisites
* A MongoDB instance (tested with MongoDB v4)
* A TLS key and certificate for providing mTLS encryption for the backend REST API (see below)
* Docker Engine (recommended)

## Generating Certificates
By default, the server will create a CA and generate keys and certificates for both server and client on startup (see ```BACKEND_GENERATE_CERT``` [configuration option](config.md)).

Alternatively, you can generate the server's and clients' keys and certificates manually using OpenSSL:

```
# Generate CA Key & Certificate
openssl genrsa -out certs/ca.key 4096
openssl req -new -x509 -sha256 -key certs/ca.key -out certs/ca.crt -days 3650

# Generate Server Key, CSR & sign with CA
openssl genrsa -out certs/server.key 4096
openssl req -new -key certs/server.key -out certs/server.csr
openssl x509 -req -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt

# Generate Application/Client Key, CSR & sign with CA
openssl genrsa -out certs/client.key 4096
openssl req -new -key certs/client.key -out certs/client.csr
openssl x509 -req -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -out certs/client.crt
```

The key and certificate must be stored in the directory specified by the ```BACKEND_CERT_DIR``` option, which is ```./certs/``` if not set (see [Configuration](config.md)). Please use these filenames:

* CA Certificate: ```ca.crt```
* Server Certificate: ```server.crt```
* Server Key: ```server.key```

## Running in Docker
It's recommended to use the pre-built Docker images to run JWT Auth Proxy. The images are built automatically with each new version and pushed to the Docker Hub. They are multi-arch, thus the correct image for your architecture will be used automatically (AMD64, ARM v6, ARM v7 and ARM64 v8).

You can use JWT Auth Proxy without Compose or Kubernetes, but these make it easier to orchestrate JWT Auth Proxy with your frontend, backend and the MongoDB.

Use the following command to run the image and set the [configuration options](config.md) accordingly:

```
docker run -d \
    -e "JWT_SIGNING_KEY=<Your JWT Signing Key>" \
    -e "MONGO_DB_URL=mongodb://localhost:27017" \
    -e "PROXY_TARGET=http://localhost:8090" \
    -e "SMTP_SERVER=localhost:25" \
    -v ${PWD}/certs:/app/certs \
    -p 8080:8080 \
    -p 8443:8443 \
    virtualzone/jwt-auth-proxy
```

## Running in Docker Compose
Please refer to the [docker-compose.yml](https://github.com/virtualzone/jwt-auth-proxy/blob/master/example/docker-compose.yml) example on how to use the pre-build Docker image with Docker Compose.

## Running in Kubernetes
You can run JWT Auth Proxy in Kubernetes. Currently, there is no Helm Chart available, this may change in the future. In the meanwhile, please set up JWT Auth Proxy as a Pod manually by using the pre-build Docker image ```virtualzone/jwt-auth-proxy```.