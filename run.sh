#!/bin/sh
cd src/
JWT_SIGNING_KEY=AmxWhqyWDpp78ZdetaqYF5qA6vVfwjAzquvFFPpHGJZAyC4y44DBVUtTrfnQe9XZkmmDJ6LmkyQttjMXMED8RA78pW6TYuck7f3SjGaDm4rchr6AKvsdzCC8Tke4wqjv \
BACKEND_CERT_DIR=/tmp/ \
TEMPLATE_SIGNUP=../res/signup.tpl \
TEMPLATE_CHANGE_EMAIL=../res/changeemail.tpl \
TEMPLATE_RESET_PASSWORD=../res/resetpassword.tpl \
TEMPLATE_NEW_PASSWORD=../res/newpassword.tpl \
PROXY_TARGET=http://localhost:8090 \
CORS_ENABLE=1 \
BACKEND_GENERATE_CERT=1 \
TOTP_ENABLE=1 \
TOTP_ENCRYPT_KEY=w66iO0l3Kru7Qgpx \
PROXY_WHITELIST=/foo/bar \
go run `ls *.go | grep -v _test.go`