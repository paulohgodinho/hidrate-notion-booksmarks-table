#!/usr/bin/env bash

set -euo pipefail

# Configuration
BINARY_NAME="hidrate-notion-bookmarks"
BUILD_DIR="bin"
CMD_PATH="./cmd/example/main.go"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Building ${BINARY_NAME}...${NC}"

# Clean previous builds
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

# Build flags for idempotent, reproducible builds
BUILD_FLAGS=(
    -trimpath
    -ldflags="-s -w -buildid="
)

# Build for macOS ARM64
echo -e "${GREEN}Building for macOS ARM64...${NC}"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
    "${BUILD_FLAGS[@]}" \
    -o "${BUILD_DIR}/${BINARY_NAME}-darwin-arm64" \
    "${CMD_PATH}"

# Build for Linux AMD64
echo -e "${GREEN}Building for Linux AMD64...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    "${BUILD_FLAGS[@]}" \
    -o "${BUILD_DIR}/${BINARY_NAME}-linux-amd64" \
    "${CMD_PATH}"

# Generate checksums
echo -e "${GREEN}Generating checksums...${NC}"
cd "${BUILD_DIR}"

# Generate SHA256 checksums
if command -v shasum &> /dev/null; then
    shasum -a 256 ${BINARY_NAME}-* > checksums.txt
elif command -v sha256sum &> /dev/null; then
    sha256sum ${BINARY_NAME}-* > checksums.txt
else
    echo "Warning: Neither shasum nor sha256sum found. Skipping checksum generation."
fi

cd ..

# Display results
echo -e "\n${GREEN}Build complete!${NC}"
echo -e "Artifacts:"
ls -lh "${BUILD_DIR}"

if [ -f "${BUILD_DIR}/checksums.txt" ]; then
    echo -e "\n${BLUE}Checksums:${NC}"
    cat "${BUILD_DIR}/checksums.txt"
fi
