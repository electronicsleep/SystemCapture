FROM debian:stable

MAINTAINER Chris Robertson https://github.com/electronicsleep

RUN mkdir -p /usr/src/app
RUN apt-get update && apt-get install golang -y

ADD SystemCapture.go /usr/src/app

WORKDIR /usr/src/app
EXPOSE 5000

CMD ["go", "run", "SystemCapture.go"]
