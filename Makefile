VERSION = $(shell git tag -l --sort=-creatordate | head -n 1)
COMMIT = $(shell git describe --always)
BUILD_DATE = $(shell date +%Y-%m-%d)
LDFLAGS = -ldflags="-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildDate=$(BUILD_DATE)'"

.PHONY: all windows macintel macarm linux rpi prepare

BINARY_NAME=dashgoat
SOURCE_FILE=./cmd/dashgoat
WEB_DIR=./web

all: prepare windows macintel macarm linux rpi

prepare:
	cp -R $(WEB_DIR) $(SOURCE_FILE)/web

windows: prepare
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME).exe $(SOURCE_FILE)

macintel: prepare
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-mac $(SOURCE_FILE)

macarm: prepare
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-mac-arm $(SOURCE_FILE)

linux: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME) $(SOURCE_FILE)

rpi: prepare
	GOOS=linux GOARCH=arm GOARM=5 go build $(LDFLAGS) -o build/$(BINARY_NAME)-rpi $(SOURCE_FILE)

ci:
	./tests/link-bin.sh
	./tests/start-single.sh
	./tests/ttl-test.sh
	./tests/nextupdate-test.sh
	./tests/metrics-test.sh
	./tests/tags-test.sh
	./tests/stop-instances.sh

clean:
	./tests/stop-instances.sh
	rm -f build/dashgoa*
	rm -rf $(SOURCE_FILE)/web
	touch dashgoat
	rm dashgoat 
