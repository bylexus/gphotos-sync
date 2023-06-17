RELEASE_DIR := ./release
BUILD_DIR := .
BIN_NAME=gphotos-sync

# default target: build for your local os
$(BIN_NAME): $(wildcard *.go **/*.go)
	go build -o $(BUILD_DIR)/$@

# build for windows
$(BIN_NAME).exe: $(wildcard *.go **/*.go)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$@

# cross-compile for several OSes
.PHONY:
release:
	GOOS=windows GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-win-amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o $(RELEASE_DIR)/$(BIN_NAME)-linux-arm64

# make everything tidy
.PHONY:
clean:
	rm -f $(BUILD_DIR)/$(BIN_NAME)
	rm -f $(BUILD_DIR)/$(BIN_NAME).exe
	rm -rf $(RELEASE_DIR)