build:
	go fmt .
	go build -o sc

check:
	go fmt .
	go test -v

test: build
	go fmt .
	./sc -t -1

install: build
	sudo cp sc /usr/local/bin
