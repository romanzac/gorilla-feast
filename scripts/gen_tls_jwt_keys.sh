#!/bin/sh

# Generate RSA keypair for TLS and JWT key for API host.

OUT_DIR=./keys
COMMON_NAME=$1
ALT_NAME=$2

# RSA keypair for TLS
openssl req -newkey rsa:2048 -keyout "${OUT_DIR}/$COMMON_NAME-private.pem" -x509 -days 365 \
     -subj "/C=TW/O=UI/CN=$COMMON_NAME" -out "${OUT_DIR}/$COMMON_NAME-public.crt" -nodes -addext "subjectAltName = DNS:$COMMON_NAME,DNS:$ALT_NAME"

# RSA keypair for JWT
openssl genrsa -out "${OUT_DIR}/$COMMON_NAME-jwt-private.pem" 2048
openssl rsa -in "${OUT_DIR}/$COMMON_NAME-jwt-private.pem" -pubout -outform PEM -out "${OUT_DIR}/$COMMON_NAME-jwt-public.pem"

