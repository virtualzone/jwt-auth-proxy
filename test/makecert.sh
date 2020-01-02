#!/bin/sh
echo "Cleaning up..."
rm -rf server/ clients/
mkdir -p server/ clients/

echo "Generating Server Key..."
openssl genrsa -out server/server.key 4096

echo "Generating Server Cert..."
openssl req -new -x509 -sha256 -key server/server.key -out server/server.crt -days 3650 \
    -subj "/C=DE/ST=/L=Frankfurt/O=Virtualzone/OU=Dev/CN=localhost"

echo "00" > server/server.srl

echo "Generating key for client 1..."
openssl genrsa -out clients/client1.key 4096
openssl req -new -key clients/client1.key -out clients/client1.csr \
    -subj "/C=DE/ST=/L=Frankfurt/O=Virtualzone/OU=Dev/CN=client1"
openssl x509 -req -in clients/client1.csr -CA server/server.crt -CAkey server/server.key -CAserial server/server.srl -out clients/client1.crt

echo "Generating key for client 2..."
openssl genrsa -out clients/client2.key 4096
openssl req -new -key clients/client2.key -out clients/client2.csr \
    -subj "/C=DE/ST=/L=Frankfurt/O=Virtualzone/OU=Dev/CN=client2"
openssl x509 -req -in clients/client2.csr -CA server/server.crt -CAkey server/server.key -CAserial server/server.srl -out clients/client2.crt
