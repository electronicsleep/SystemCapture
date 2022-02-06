#!/bin/bash
set +e
docker rm systemcapture
set -e
GOOS=linux go build SystemCapture.go
docker build -t systemcapture .
docker run -p 8080:8080 --name systemcapture -i -t systemcapture
