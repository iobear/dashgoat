.PHONY: all windows macintel macarm linux rpi prepare

BINARY_NAME=dashgoat
SOURCE_FILE=./cmd/dashgoat
WEB_DIR=./web

all: prepare windows macintel macarm linux rpi

prepare:
	cp -R $(WEB_DIR) $(SOURCE_FILE)/web

windows: prepare
	GOOS=windows GOARCH=amd64 go build -o build/$(BINARY_NAME).exe $(SOURCE_FILE)

macintel: prepare
	GOOS=darwin GOARCH=amd64 go build -o build/$(BINARY_NAME)-mac $(SOURCE_FILE)

macarm: prepare
	GOOS=darwin GOARCH=arm64 go build -o build/$(BINARY_NAME)-mac-arm $(SOURCE_FILE)

linux: prepare
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) $(SOURCE_FILE)

rpi: prepare
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/$(BINARY_NAME)-rpi $(SOURCE_FILE)