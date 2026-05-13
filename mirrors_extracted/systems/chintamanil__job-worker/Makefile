# Job Worker Service Makefile
#
# Targets:
#   make              Build server and client
#   make test         Run tests (local)
#   make test-docker  Run tests in Docker (Linux)
#   make test-docker-privileged  Run tests with cgroups (Linux, privileged)
#   make proto        Generate protobuf code
#   make certs        Generate TLS certificates
#   make clean        Remove build artifacts
#   make help         Show all targets
#
# OrbStack Integration Tests (macOS):
#   make setup-orb           Setup OrbStack VM (run once)
#   make integration-check-orb   Verify OrbStack VM is ready
#   make test-integration-orb    Run integration tests in OrbStack VM

BUILDDIR := build
UNAME_S := $(shell uname -s)
DOCKER_IMAGE := job-worker-test
ORB_VM := jobworker

.PHONY: all build server server-linux client proto test test-race test-docker test-docker-privileged docker-build lint certs clean help
.PHONY: setup-orb integration-check-orb test-integration-orb

all: build

build: server client

server: | $(BUILDDIR)
ifeq ($(UNAME_S),Darwin)
	@echo "Note: Server requires Linux (cgroups v2)"
endif
	go build -o $(BUILDDIR)/worker-server ./cmd/worker-server

client: | $(BUILDDIR)
	go build -o $(BUILDDIR)/worker-cli ./cmd/worker-cli

# Build server for Linux (for running in OrbStack VM)
server-linux: | $(BUILDDIR)
	orb -m $(ORB_VM) /usr/local/go/bin/go build -C $(CURDIR) -o $(BUILDDIR)/worker-server-linux ./cmd/worker-server

proto:
	@mkdir -p pkg/proto
	protoc --proto_path=api/proto \
		--go_out=pkg/proto --go_opt=paths=source_relative \
		--go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
		api/proto/worker.proto

test:
	go test -v -cover ./...

test-race:
	go test -race -v ./...

# Build Docker image for testing
docker-build:
	docker build -f Dockerfile.testing -t $(DOCKER_IMAGE) .

# Run tests in Docker (Linux environment, no privileges)
test-docker: docker-build
	docker run --rm $(DOCKER_IMAGE)

# Run tests in Docker with cgroup access (for integration tests)
test-docker-privileged: docker-build
	docker run --rm --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:rw $(DOCKER_IMAGE)

lint:
	go fmt ./...
	go vet ./...

certs:
	@cd certs && $(MAKE) all

clean:
	rm -rf $(BUILDDIR)
	-docker rmi $(DOCKER_IMAGE) 2>/dev/null

$(BUILDDIR):
	@mkdir -p $(BUILDDIR)

# ============================================================================
# OrbStack Integration Tests (macOS)
# ============================================================================

# Setup OrbStack VM with all dependencies (run once)
setup-orb:
	@./scripts/setup-orbstack.sh

# Verify OrbStack VM is ready for integration tests
integration-check-orb:
	@echo "==> Checking OrbStack VM '$(ORB_VM)'..."
	@orb list | grep -q "^$(ORB_VM) " || (echo "ERROR: VM '$(ORB_VM)' not found. Run 'make setup-orb' first." && exit 1)
	@echo "==> Checking cgroups v2..."
	@orb -m $(ORB_VM) cat /sys/fs/cgroup/cgroup.controllers
	@echo "==> Checking stress-ng..."
	@orb -m $(ORB_VM) which stress-ng
	@echo "==> Checking Go..."
	@orb -m $(ORB_VM) /usr/local/go/bin/go version
	@echo ""
	@echo "VM '$(ORB_VM)' is ready for integration tests."

# Run integration tests in OrbStack VM
# Uses -m to specify machine, runs as root for cgroup access
# Note: -race disabled because it requires CGO which adds complexity
test-integration-orb: integration-check-orb
	@echo "==> Running integration tests in OrbStack VM..."
	orb -m $(ORB_VM) -u root /usr/local/go/bin/go test -C $(CURDIR) -v -tags=integration -count=1 ./test/integration/...

# ============================================================================

help:
	@echo "Build:    make, make server, make client, make server-linux"
	@echo "Test:     make test, make test-race"
	@echo "Docker:   make test-docker, make test-docker-privileged"
	@echo "OrbStack: make setup-orb, make integration-check-orb, make test-integration-orb"
	@echo "Other:    make proto, make certs, make lint, make clean"
