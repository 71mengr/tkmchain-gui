BINARY_NAME=gtkm-wallet
VERSION=1.0.0
BUILD_DIR=build
GO=go

.PHONY: all clean build-linux build-windows build-mac install-deps

all: build-linux build-windows build-mac

build-linux:
	@echo "Building for Linux..."
	mkdir -p $(BUILD_DIR)/linux
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO) build -tags "production" -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/linux/$(BINARY_NAME) ./cmd/gtkm-wallet
	@echo "✅ Linux binary: $(BUILD_DIR)/linux/$(BINARY_NAME)"

build-windows:
	@echo "Building for Windows..."
	mkdir -p $(BUILD_DIR)/windows
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(GO) build -tags "production" -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/windows/$(BINARY_NAME).exe ./cmd/gtkm-wallet
	@echo "✅ Windows binary: $(BUILD_DIR)/windows/$(BINARY_NAME).exe"

build-windows-cross:
	@echo "Building for Windows (cross-compile)..."
	mkdir -p $(BUILD_DIR)/windows
	# For cross-compilation, we need to use the go-win bindings
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc $(GO) build -tags "production windows" -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/windows/$(BINARY_NAME).exe ./cmd/gtkm-wallet
	@echo "✅ Windows binary: $(BUILD_DIR)/windows/$(BINARY_NAME).exe"

build-mac:
	@echo "Building for macOS..."
	mkdir -p $(BUILD_DIR)/mac
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GO) build -tags "production" -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/mac/$(BINARY_NAME) ./cmd/gtkm-wallet
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(GO) build -tags "production" -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/mac/$(BINARY_NAME)-arm64 ./cmd/gtkm-wallet
	@echo "✅ macOS binaries: $(BUILD_DIR)/mac/$(BINARY_NAME) and $(BUILD_DIR)/mac/$(BINARY_NAME)-arm64"

build-all: build-linux build-windows build-mac

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

install-deps:
	$(GO) mod tidy
	$(GO) get fyne.io/fyne/v2@latest
	$(GO) get github.com/ethereum/go-ethereum@latest
	
	# Install Windows cross-compilation tools if needed
	@echo "For Windows cross-compilation, install:"
	@echo "  sudo apt-get install gcc-mingw-w64-x86-64"

run:
	CGO_ENABLED=1 $(GO) run ./cmd/gtkm-wallet

test:
	CGO_ENABLED=1 $(GO) test ./...

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build-linux      - Build Linux binary"
	@echo "  make build-windows    - Build Windows binary (requires mingw)"
	@echo "  make build-mac        - Build macOS binaries"
	@echo "  make build-all        - Build all platforms"
	@echo "  make run              - Run the wallet"
	@echo "  make clean            - Clean build directory"
	@echo "  make install-deps     - Install dependencies"
	@echo "  make test             - Run tests"
