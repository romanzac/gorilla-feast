#!/bin/sh
# Start script for Gorilla Feast API

# Copy and import certificates. Can be potentially added later externally.
cp -Rp /app/*.crt /usr/local/share/ca-certificates
update-ca-certificates

# Start API service process
./gorilla-feast-linux

