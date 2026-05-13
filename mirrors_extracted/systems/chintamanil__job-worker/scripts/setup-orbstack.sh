#!/bin/bash
# Setup OrbStack VM for integration tests
#
# This script creates and configures an Ubuntu VM with all dependencies
# needed to run the job-worker integration tests.
#
# Prerequisites:
#   - OrbStack installed (brew install orbstack)
#
# Usage:
#   ./scripts/setup-orbstack.sh

set -euo pipefail

VM_NAME="jobworker"
PROJECT_DIR="/Users/manil/Projects/Go/job-worker"

echo "==> Checking OrbStack installation..."
if ! command -v orb &> /dev/null; then
    echo "ERROR: OrbStack not installed. Install with: brew install orbstack"
    exit 1
fi

echo "==> Checking if VM '$VM_NAME' exists..."
if orb list | grep -q "^$VM_NAME "; then
    echo "VM '$VM_NAME' already exists."
    read -p "Delete and recreate? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "==> Deleting existing VM..."
        orb delete "$VM_NAME"
    else
        echo "Keeping existing VM. Run 'make integration-check-orb' to verify."
        exit 0
    fi
fi

echo "==> Creating Ubuntu VM '$VM_NAME'..."
orb create ubuntu "$VM_NAME"

echo "==> Waiting for VM to be ready..."
sleep 3

echo "==> Updating package lists..."
orb -m "$VM_NAME" sudo apt-get update -qq

echo "==> Installing dependencies..."
orb -m "$VM_NAME" sudo apt-get install -y -qq \
    stress-ng \
    build-essential \
    curl \
    > /dev/null

echo "==> Installing Go..."
orb -m "$VM_NAME" bash -c '
    if [ ! -f /usr/local/go/bin/go ]; then
        echo "Downloading Go..."
        curl -sLO https://go.dev/dl/go1.23.6.linux-arm64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go1.23.6.linux-arm64.tar.gz
        rm go1.23.6.linux-arm64.tar.gz
    fi
    /usr/local/go/bin/go version
'

echo "==> Enabling cgroup controllers..."
# Enable cpu, memory, io controllers at root level
orb -m "$VM_NAME" sudo bash -c 'echo "+cpu +memory +io" > /sys/fs/cgroup/cgroup.subtree_control' || true

# Create test cgroup directory and enable controllers for it
orb -m "$VM_NAME" sudo bash -c '
    mkdir -p /sys/fs/cgroup/jobworker-integration-test
    echo "+cpu +memory +io" > /sys/fs/cgroup/jobworker-integration-test/cgroup.subtree_control
' || true

echo "==> Verifying cgroups v2..."
CONTROLLERS=$(orb -m "$VM_NAME" cat /sys/fs/cgroup/cgroup.controllers 2>/dev/null || echo "")
if [[ -z "$CONTROLLERS" ]]; then
    echo "ERROR: cgroups v2 not available in VM"
    exit 1
fi
echo "Available controllers: $CONTROLLERS"

echo "==> Verifying subtree controllers..."
orb -m "$VM_NAME" cat /sys/fs/cgroup/cgroup.subtree_control || true

echo "==> Verifying stress-ng..."
orb -m "$VM_NAME" stress-ng --version | head -1

echo "==> Verifying Go..."
orb -m "$VM_NAME" /usr/local/go/bin/go version

echo "==> Building Linux server binary..."
orb -m "$VM_NAME" /usr/local/go/bin/go build -C "$PROJECT_DIR" -o "$PROJECT_DIR/build/worker-server-linux" ./cmd/worker-server

echo ""
echo "==> Setup complete!"
echo ""
echo "VM '$VM_NAME' is ready for integration tests."
echo ""
echo "Run tests with:"
echo "  make test-integration-orb"
echo ""
echo "Check VM status with:"
echo "  make integration-check-orb"
echo ""
echo "Start server manually with:"
echo "  cd $PROJECT_DIR && orb -m $VM_NAME -u root ./build/worker-server-linux"
echo ""
echo "Run CLI commands with:"
echo "  orb -m $VM_NAME /usr/local/go/bin/go run -C $PROJECT_DIR ./cmd/worker-cli start /bin/echo hello"
