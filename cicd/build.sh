#!/bin/bash
mkdir bin
set -ex
export PATH=$PATH:/usr/local/go/bin
go version
echo "Building for Linux and MacOS"
GOOS=linux go build -o bin/SystemCapture SystemCapture.go
GOOS=darwin go build -o bin/SystemCapture-MacOS SystemCapture.go
echo "Done"
