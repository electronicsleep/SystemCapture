FROM ubuntu:22.04

LABEL org.opencontainers.image.authors="https://github.com/electronicsleep"

RUN mkdir -p /usr/src/app

RUN apt-get update && apt-get upgrade -y && apt-get install -y net-tools

ADD SystemCapture /usr/src/app

WORKDIR /usr/src/app
EXPOSE 8080

# Run Webserver mode
# CMD ["./SystemCapture", "-t", "1", "-w"]
# Run Console mode
CMD ["./SystemCapture", "-t", "-1"]
