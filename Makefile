GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

BINARY_NAME=gin_simple_api
MAIN_PATH=$(wildcard *.go)

mod:
	make go.mod
	make go.sum

go.mod:
	$(GOMOD) init gin_simple_api

go.sum:
	$(GOGET) -d -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f go.mod go.sum *.db

run: mod
	go run $(MAIN_PATH)
