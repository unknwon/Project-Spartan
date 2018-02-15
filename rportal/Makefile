.PHONY: build run

all: build run

build:
	go install -v
	cp '$(GOPATH)/bin/rportal' .

run:
	./rportal -name=$(name) -port=$(port)