# Project variables
PROJECT_NAME := docker-machine-driver-xelon

# Build variables
BUILD_DIR := build
DEV_GOARCH := $(shell go env GOARCH)
DEV_GOOS := $(shell go env GOOS)


## clean: Delete the build directory.
.PHONY: clean
clean:
	@echo "==> Removing '$(BUILD_DIR)' directory..."
	@rm -rf $(BUILD_DIR)


## fmt: Format code with go fmt.
.PHONY: fmt
fmt:
	@echo "==> Checking code with 'go fmt'..."
	@go fmt ./...


## test: Run all tests.
.PHONY: test
test:
	@echo "==> Running tests..."
	@mkdir -p $(BUILD_DIR)
	@go test -v -cover -coverprofile=$(BUILD_DIR)/coverage.out ./...


## build: Build binary for default local system's operating system and architecture.
.PHONY: build
build:
	@echo "==> Building binary..."
	@echo "    running go build for GOOS=$(DEV_GOOS) GOARCH=$(DEV_GOARCH)"
ifeq ($(OS),Windows_NT)
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME)_$(DEV_GOOS)_$(DEV_GOARCH).exe cmd/main.go
else
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME)_$(DEV_GOOS)_$(DEV_GOARCH) cmd/main.go
endif


## release: Build release binaries for all supported versions.
.PHONY: release
release:
	@echo "==> Building release binaries..."
	@echo "    running go build for GOOS=darwin GOARCH=amd64"
	@GOARCH=amd64 GOOS=darwin go build -o $(BUILD_DIR)/$(PROJECT_NAME)_darwin_amd64 cmd/main.go
	@echo "    running go build for GOOS=linux GOARCH=amd64"
	@GOARCH=amd64 GOOS=linux go build -o $(BUILD_DIR)/$(PROJECT_NAME)_linux_amd64 cmd/main.go
	@echo "    running go build for GOOS=windows GOARCH=amd64"
	@GOARCH=amd64 GOOS=windows go build -o $(BUILD_DIR)/$(PROJECT_NAME)_windows_amd64.exe cmd/main.go
	@echo "==> Generate checksums..."
	@cd $(BUILD_DIR) && for f in *; do sha256sum "$$f" > "$$f.sha256"; done


help: Makefile
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
