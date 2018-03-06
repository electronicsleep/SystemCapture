#!/bin/bash
GOOS=linux go build SystemCapture.go
docker rm scapp
docker build -t scapp .
docker run -p 5000:5000 --name scapp -i -t scapp
