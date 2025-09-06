.PHONY: build clean test run

# Variables
BINARY_NAME_CLI=forgetti-cli
BINARY_NAME_SERVER=forgetti-server
BINARY_DIR=bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags="-X 'Forgetti/cmd.Version=$(VERSION)' -X 'Forgetti/cmd.Commit=$(COMMIT)' -X 'Forgetti/cmd.BuildDate=$(BUILD_DATE)'"

# Tasks
build:
	make build-cli
	make build-server

build-cli:
	mkdir -p $(BINARY_DIR)
	cd cli_tool && go build $(LDFLAGS) -o ../$(BINARY_DIR)/$(BINARY_NAME_CLI)
	cp cli_tool/config.json $(BINARY_DIR)/.config.json

build-server:
	mkdir -p $(BINARY_DIR)
	cd server && go build $(LDFLAGS) -o ../$(BINARY_DIR)/$(BINARY_NAME_SERVER)
	cp server/config.json $(BINARY_DIR)/config.json

clean:
	rm -rf $(BINARY_DIR)
	go clean

test:
	make test-common
	make test-cli
	make test-server

test-common:
	cd common && go test ./... -v

test-cli:
	cd cli_tool && go test ./... -v
	
test-server:
	cd server && go test ./... -v

deps:
	make deps-cli
	make deps-server

deps-cli:
	cd cli_tool && go mod download

deps-server:
	cd server && go mod download

