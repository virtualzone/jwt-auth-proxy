#!/bin/sh
curl -v \
    --cacert server/server.crt \
    --key clients/client1.key \
    --cert clients/client1.crt \
    https://localhost:8443