#!/bin/bash

# Build gorilla-feast for Linux AMD64
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o gorilla-feast-linux

docker build --no-cache -t "gorilla-feast:1.0.2" -f Dockerfile .

