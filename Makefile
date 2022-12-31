build:
	go fmt .; go build -o sc

run:
	./sc -t -1

install:
	sudo cp sc /usr/local/bin
