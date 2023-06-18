RELEASE_DIR := ./release
BUILD_DIR := .
BIN_NAME = gphotos-sync

.PHONY: release clean

# default target: build for your local os
$(BIN_NAME): $(wildcard *.go **/*.go) go.mod go.sum
	go build -o $(BUILD_DIR)/$@

# build for windows
$(BIN_NAME).exe: $(wildcard *.go **/*.go) go.mod go.sum
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$@

# cross-compile for several OSes
release:
	TAG=$(shell git describe --tag); \
	GOOS=windows GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-win-amd64-$${TAG}.exe; \
	GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-darwin-amd64-$${TAG}; \
	GOOS=darwin GOARCH=arm64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-darwin-arm64-$${TAG}; \
	GOOS=linux GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-linux-amd64-$${TAG}; \
	GOOS=linux GOARCH=arm64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-linux-arm64-$${TAG};

# make everything tidy
clean:
	rm -f $(BUILD_DIR)/$(BIN_NAME); \
	rm -f $(BUILD_DIR)/$(BIN_NAME).exe; \
	rm -rf $(RELEASE_DIR);