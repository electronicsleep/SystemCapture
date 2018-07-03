FROM alpine:latest

MAINTAINER Chris Robertson https://github.com/electronicsleep

RUN mkdir -p /usr/src/app

#For verbose commands
RUN apk update && apk add sysstat net-tools lsof procps

ADD SystemCapture /usr/src/app

WORKDIR /usr/src/app
EXPOSE 5000

#Run normally
CMD ["./SystemCapture", "-v"]
