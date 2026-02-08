#!/usr/bin/env bash

set -euo pipefail

BINARY_NAME="hidrate-notion-bookmarks"
BUILD_DIR="bin"

# Build flags to remove all VCS info and metadata
# -trimpath removes absolute paths
# -ldflags="-s -w" removes symbols and debug info
# -buildvcs=false prevents embedding VCS info (go 1.18+)
GO_BUILD_FLAGS="-trimpath -ldflags=-s -ldflags=-w -buildvcs=false"

# Clean and create bin directory
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

cd src

# Build for Linux AMD64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GO_BUILD_FLAGS} -o "../${BUILD_DIR}/${BINARY_NAME}-linux-amd64" .

# Build for macOS ARM64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${GO_BUILD_FLAGS} -o "../${BUILD_DIR}/${BINARY_NAME}-darwin-arm64" .

# Build for Windows AMD64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${GO_BUILD_FLAGS} -o "../${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe" .

cd ..

echo "✓ Binaries built in ${BUILD_DIR}/"
