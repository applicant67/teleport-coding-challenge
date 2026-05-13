.PHONY: all build build-linux test test-race test-cover test-cover-html clean proto tools check certs

# Output directories and names
BIN_DIR := bin
COVERAGE_DIR := coverage
SERVER_BIN := worker-server
CLI_BIN := worker-cli

# Protobuf settings
PROTO_DIR := api/v1

# Build parameters
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
EXT :=
ifeq ($(GOOS),windows)
	EXT := .exe
endif

# Race detection requires a C compiler (GCC).
# On Windows, we check if it's available to avoid build failures.
RACE_FLAG := -race
ifeq ($(GOOS),windows)
	ifeq ($(shell where gcc 2> NUL),)
		RACE_FLAG :=
	endif
endif

# Export environment variables so they are passed to child processes
export GOOS
export GOARCH
export CGO_ENABLED

all: build

build:
	@echo "Building binaries ($(GOOS)/$(GOARCH))..."
ifeq ($(GOOS),linux)
	go build -buildvcs=false -o $(BIN_DIR)/$(SERVER_BIN)$(EXT) ./cmd/server
endif
	go build -buildvcs=false -o $(BIN_DIR)/$(CLI_BIN)$(EXT) ./cmd/cli

build-linux:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-windows:
	@$(MAKE) build GOOS=windows GOARCH=amd64

test:
	@echo "Running tests..."
	go test -buildvcs=false -race $$(go list -buildvcs=false ./... | grep -v /api/v1)

test-root:
	@echo "Running tests with sudo..."
	go test -buildvcs=false -exec sudo $$(go list -buildvcs=false ./... | grep -v /api/v1)

test-race:
	@echo "Running tests..."
ifneq ($(RACE_FLAG),)
	@echo "Race detection enabled."
	go test -buildvcs=false $(RACE_FLAG) $$(go list -buildvcs=false ./... | grep -v /api/v1)
else
	@echo "WARNING: GCC not found. Skipping race detection."
	go test -buildvcs=false $$(go list -buildvcs=false ./... | grep -v /api/v1)
endif

test-cover:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	go test -buildvcs=false -coverprofile=$(COVERAGE_DIR)/coverage.out $$(go list -buildvcs=false ./... | grep -v /api/v1)
	go tool cover -func=$(COVERAGE_DIR)/coverage.out
ifdef SUDO_USER
	@chown $(SUDO_USER) $(COVERAGE_DIR)/coverage.out
endif

test-cover-html: test-cover
	@echo "Generating coverage report..."
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
ifdef SUDO_USER
	@chown $(SUDO_USER) $(COVERAGE_DIR)/coverage.html
endif

test-complete:
	@tr -d '\r' < scripts/test-complete.sh | bash -s -- $(COVERAGE_DIR)

clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR) $(COVERAGE_DIR)

proto:
	@echo "Generating protobuf code..."
	protoc --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/worker.proto

tools:
	@echo "Installing Go protobuf plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

check:
	@echo "Running formatting..."
	go fmt ./...
	@echo "Running static analysis..."
	go vet ./...
	@$(MAKE) test-race CGO_ENABLED=1
	@echo "Tidying dependencies..."
	go mod tidy

certs: ## Generate mTLS certificates for development and testing
	@tr -d '\r' < scripts/generate_certs.sh > scripts/generate_certs.tmp.sh
	@bash scripts/generate_certs.tmp.sh
	@rm scripts/generate_certs.tmp.sh
