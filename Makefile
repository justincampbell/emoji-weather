HOMEPAGE=https://github.com/justincampbell/emoji-weather
PREFIX=/usr/local

COVERAGE_FILE = coverage.out

ARCHIVE=emoji-weather-$(TAG).tar.gz
ARCHIVE_URL=$(HOMEPAGE)/archive/$(TAG).tar.gz

test: acceptance

install: build
	mkdir -p $(PREFIX)/bin
	cp -v bin/emoji-weather $(PREFIX)/bin/emoji-weather

uninstall:
	rm -vf $(PREFIX)/bin/emoji-weather

coverage: unit
	go tool cover -html=$(COVERAGE_FILE)

acceptance: build
	bats test

build: dependencies unit
	go build -o bin/emoji-weather

unit: dependencies
	go test -coverprofile=$(COVERAGE_FILE) -timeout 25ms

dependencies:
	go get -t
	go get golang.org/x/tools/cmd/cover

.PHONY: acceptance build coverage dependencies install test uninstall unit
