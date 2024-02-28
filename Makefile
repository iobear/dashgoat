.PHONY: all windows macintel macarm linux rpi

BINARY_NAME=dashgoat
SOURCE_FILE=./cmd/dashgoat

all: windows macintel macarm linux rpi

windows:
	GOOS=windows GOARCH=amd64 go build -o build/$(BINARY_NAME).exe $(SOURCE_FILE)

macintel:
	GOOS=darwin GOARCH=amd64 go build -o build/$(BINARY_NAME)-mac $(SOURCE_FILE)

macarm:
	GOOS=darwin GOARCH=arm64 go build -o build/$(BINARY_NAME)-mac-arm $(SOURCE_FILE)

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) $(SOURCE_FILE)

rpi:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/$(BINARY_NAME)-rpi $(SOURCE_FILE)
