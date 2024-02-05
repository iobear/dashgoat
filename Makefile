.PHONY: all windows macintel macarm linux rpi

BINARY_NAME=dashGoat
SOURCE_FILE=./cmd/dashgoat

all: windows macintel macarm linux rpi

windows:
	GOOS=windows GOARCH=amd64 go build -o build/$(BINARY_NAME)-amd64.exe $(SOURCE_FILE)

macintel:
	GOOS=darwin GOARCH=amd64 go build -o build/$(BINARY_NAME)-macintel $(SOURCE_FILE)

macarm:
	GOOS=darwin GOARCH=arm64 go build -o build/$(BINARY_NAME)-macarm $(SOURCE_FILE)

linux:
	GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) $(SOURCE_FILE)

rpi:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/$(BINARY_NAME) $(SOURCE_FILE)