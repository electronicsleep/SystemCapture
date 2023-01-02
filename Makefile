build:
	go fmt .; go build -o sc

test: build
	./sc -t -1

install:
	sudo cp sc /usr/local/bin
