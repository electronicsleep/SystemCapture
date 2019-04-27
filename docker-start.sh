#!/bin/bash
GOOS=linux go build SystemCapture.go
docker rm systemcapture
docker build -t systemcapture .
docker run -p 8080:8080 --name systemcapture -i -t systemcapture
