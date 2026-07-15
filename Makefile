VERSION = $(shell git tag --sort=-v:refname | head -n 1)
COMMIT = $(shell git describe --always)
BUILD_DATE = $(shell date +%Y-%m-%d)
LDFLAGS = -ldflags="-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildDate=$(BUILD_DATE)'"
DOCKERFILE_PATH = "build/package/Dockerfile"

.PHONY: all windows macintel macarm linux rpi prepare

BINARY_NAME=dashgoat
SOURCE_FILE=./cmd/dashgoat
WEB_DIR=./web

all: prepare windows macintel macarm linux rpi
docker: podman

prepare:
	cp -R $(WEB_DIR) $(SOURCE_FILE)/web

build: prepare
	CGO_ENABLED=0 go build $(LDFLAGS) -o . $(SOURCE_FILE)

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

podman: prepare
	podman build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg BUILD_DATE=$(BUILD_DATE)  -f $(DOCKERFILE_PATH) -t docker.io/analogbear/dashgoat:$(VERSION) -t docker.io/analogbear/dashgoat:latest .

test: build
	./tests/link-bin.sh
	TEST_LOGFILE=/tmp/ttl-test.log ./tests/start-single.sh && ./tests/ttl-test.sh && ./tests/stop-instances.sh
	TEST_LOGFILE=/tmp/nextupdate-test.log ./tests/start-single.sh && ./tests/nextupdate-test.sh && ./tests/stop-instances.sh
	TEST_LOGFILE=/tmp/metrics-test.log ./tests/start-single.sh && ./tests/metrics-test.sh && ./tests/stop-instances.sh
	TEST_LOGFILE=/tmp/tags-test.log ./tests/start-single.sh && ./tests/tags-test.sh && ./tests/stop-instances.sh
	TEST_LOGFILE=/tmp/heartbeat-test.log ./tests/start-single.sh && ./tests/heartbeat-test.sh && ./tests/stop-instances.sh
	TEST_LOGFILE=/tmp/alertmanager-test.log ./tests/start-single.sh && ./tests/alertmanager-test.sh && ./tests/stop-instances.sh
	TEST_LOGFILE=/tmp/search-test.log ./tests/start-single.sh && ./tests/search-test.sh && ./tests/stop-instances.sh
	TEST_LOGFILE=/tmp/dependon-test.log ./tests/start-single.sh && ./tests/dependon-test.sh && ./tests/stop-instances.sh

clean:
	./tests/stop-instances.sh
	rm -f build/dashgoa*
	rm -rf $(SOURCE_FILE)/web
	touch dashgoat
	rm dashgoat 
