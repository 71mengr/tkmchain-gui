BINARY_NAME ?= gtkm-wallet
BUILD_DIR ?= build
GO ?= go

# Versioning: default to latest tag, fallback to 1.0.0
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo 1.0.0)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS ?= -s -w -X 'main.Version=$(VERSION)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'
BUILD_FLAGS ?= -trimpath -tags "production" -ldflags="$(LDFLAGS)"

# Default: keep CGO enabled because desktop UI libraries (fyne) require it on some platforms.
CGO_ENABLED ?= 1

.PHONY: all clean build-linux build-windows build-mac build-all build-release install-deps run test fmt vet help

all: build-all

# Linux (amd64)
build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)/linux
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/linux/$(BINARY_NAME) ./cmd/gtkm-wallet
	@echo "✅ Linux binary: $(BUILD_DIR)/linux/$(BINARY_NAME)"
	@{ command -v strip >/dev/null 2>&1 && strip $(BUILD_DIR)/linux/$(BINARY_NAME) || true; }

# Windows (amd64) - produces .exe
build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BUILD_DIR)/windows
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/windows/$(BINARY_NAME).exe ./cmd/gtkm-wallet
	@echo "✅ Windows binary: $(BUILD_DIR)/windows/$(BINARY_NAME).exe"

# macOS - build both amd64 and arm64
build-mac:
	@echo "Building for macOS (amd64 + arm64)..."
	@mkdir -p $(BUILD_DIR)/mac
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/mac/$(BINARY_NAME)-amd64 ./cmd/gtkm-wallet
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=arm64 $(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/mac/$(BINARY_NAME)-arm64 ./cmd/gtkm-wallet
	@echo "✅ macOS binaries: $(BUILD_DIR)/mac/$(BINARY_NAME)-amd64 and $(BUILD_DIR)/mac/$(BINARY_NAME)-arm64"

build-all: build-linux build-windows build-mac

# Package artifacts for release and generate checksums
build-release: build-all
	@echo "Packaging release artifacts for v$(VERSION)..."
	@rm -rf $(BUILD_DIR)/dist
	@mkdir -p $(BUILD_DIR)/dist
	# Linux tar
	cd $(BUILD_DIR)/linux && tar -czf ../dist/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)
	# Windows zip
	cd $(BUILD_DIR)/windows && zip -r ../dist/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME).exe >/dev/null
	# mac tar
	cd $(BUILD_DIR)/mac && tar -czf ../dist/$(BINARY_NAME)-$(VERSION)-darwin-universal.tar.gz $(BINARY_NAME)-amd64 $(BINARY_NAME)-arm64
	# Checksums
	cd $(BUILD_DIR)/dist && sha256sum * > $(BINARY_NAME)-$(VERSION)-SHA256SUMS
	@echo "✅ Release artifacts in $(BUILD_DIR)/dist"

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

# Install or verify development / packaging tools. Do NOT modify project modules.
install-deps:
	@echo "Tidying modules and downloading dependencies..."
	$(GO) mod tidy
	$(GO) mod download
	@echo "Installing recommended CLI tools for building/packaging (local GOPATH/bin should be on PATH)"
	@echo " - fyne build tools (for packaging desktop apps)"
	$(GO) install fyne.io/fyne/v2/cmd/fyne@latest
	@echo " - goreleaser (optional): used to produce reproducible releases"
	$(GO) install github.com/goreleaser/goreleaser@latest
	@echo "If you need to cross-compile on Debian/Ubuntu, install mingw-w64:"
	@echo "  sudo apt-get install -y gcc-mingw-w64-x86-64 mingw-w64"

# Run the app locally (development)
run:
	@echo "Running (development)..."
	CGO_ENABLED=$(CGO_ENABLED) $(GO) run -tags "production" ./cmd/gtkm-wallet

# Run tests
test:
	@echo "Running tests..."
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test ./... -v

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

vet:
	@echo "Running go vet..."
	$(GO) vet ./...

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build-linux      - Build Linux binary (amd64)"
	@echo "  make build-windows    - Build Windows binary (amd64)"
	@echo "  make build-mac        - Build macOS binaries (amd64 & arm64)"
	@echo "  make build-all        - Build all platforms"
	@echo "  make build-release    - Build and package release artifacts (tar/zip + SHA256)"
	@echo "  make run              - Run the wallet (development)"
	@echo "  make test             - Run unit tests"
	@echo "  make install-deps     - Install CLI helper tools (fyne, goreleaser)"
	@echo "  make fmt              - Run go fmt"
	@echo "  make vet              - Run go vet"
	@echo "  make clean            - Remove build artifacts"
