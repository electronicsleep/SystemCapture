FROM debian:stable

MAINTAINER Chris Robertson https://github.com/electronicsleep

RUN mkdir -p /usr/src/app

#For verbose commands
#RUN apt-get update && apt-get install sysstat net-tools lsof -y

#Test build/run go program Linux
#RUN apt-get update && apt-get install golang -y
#ADD SystemCapture.go /usr/src/app

#Test run go program
RUN apt-get update &&  apt-get install procps -y
ADD SystemCapture /usr/src/app

WORKDIR /usr/src/app

#Test build/run go program
#RUN go build SystemCapture.go

#Run normally
CMD ["./SystemCapture", "-t"]
