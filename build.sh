#!/bin/bash
mkdir bin
set -ex
echo "Building for Linux and MacOS"
gofmt -w SystemCapture.go
GOOS=linux go build SystemCapture.go
cp SystemCapture bin/SystemCapture-Linux
GOOS=darwin go build SystemCapture.go
cp SystemCapture bin/SystemCapture-MacOS
echo "Done"
