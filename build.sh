#!/bin/bash
mkdir bin
set -e
echo "Building for Linux and MacOS"
gofmt -w SystemCapture.go
GOOS=linux go build SystemCapture.go
mv SystemCapture bin/SystemCapture-Linux
GOOS=darwin go build SystemCapture.go
mv SystemCapture bin/SystemCapture-MacOS
echo "Done"
