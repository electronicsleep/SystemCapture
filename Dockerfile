FROM alpine:latest

MAINTAINER Chris Robertson https://github.com/electronicsleep

RUN mkdir -p /usr/src/app

#For verbose commands
#RUN apk update && apk install sysstat net-tools lsof

#Test run go program - install w top
#RUN apk update &&  apk install procps
ADD SystemCapture /usr/src/app

WORKDIR /usr/src/app
EXPOSE 5000

#Run normally
CMD ["./SystemCapture", "-t"]
