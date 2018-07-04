#!/bin/bash
GOOS=linux go build SystemCapture.go
docker rm systemcapture
docker build -t systemcapture .
docker run -p 5000:5000 --name systemcapture -i -t systemcapture
