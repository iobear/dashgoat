.PHONY: all windows darwin-amd64 darwin-arm64 linux

BINARY_NAME=dashGoat
SOURCE_FILE=./cmd/dashgoat

all: windows darwin-amd64 darwin-arm64 linux

windows:
	GOOS=windows GOARCH=amd64 go build -o build/$(BINARY_NAME)-amd64.exe $(SOURCE_FILE)

intelmac:
	GOOS=darwin GOARCH=amd64 go build -o build/$(BINARY_NAME)-intelmac $(SOURCE_FILE)

armmac:
	GOOS=darwin GOARCH=arm64 go build -o build/$(BINARY_NAME)-armmac $(SOURCE_FILE)

linux:
	GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) $(SOURCE_FILE)
