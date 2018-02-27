.PHONY: build run

all: build run

build:
	go install -v
	cp '$(GOPATH)/bin/haproxy' .

run:
	./haproxy -name=$(name) -port=$(port)