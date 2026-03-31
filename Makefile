.PHONY: build test install clean

BINARY  := emoji-weather
BIN_DIR := bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	mkdir -p $(BIN_DIR)
	go build -ldflags="-X main.version=$(VERSION)" -o $(BIN_DIR)/$(BINARY) .

test:
	go test ./...

install: build
	cp $(BIN_DIR)/$(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -rf $(BIN_DIR)
