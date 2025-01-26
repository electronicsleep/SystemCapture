FROM ubuntu:latest

LABEL org.opencontainers.image.authors="https://github.com/electronicsleep"

RUN mkdir -p /usr/src/app

# For verbose commands
RUN apk update && apk upgrade && apk add sysstat net-tools lsof procps

ADD SystemCapture /usr/src/app

WORKDIR /usr/src/app
EXPOSE 8080

# Run Webserver mode
# CMD ["./SystemCapture", "-t", "1", "-w"]
# Run Console mode
CMD ["./SystemCapture", "-t", "-1"]
