#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

$BINARY_NAME = "hidrate-notion-bookmarks"
$BUILD_DIR = "bin"

# Build flags to remove all VCS info and metadata
# -trimpath removes absolute paths
# -ldflags="-s -w" removes symbols and debug info
# -buildvcs=false prevents embedding VCS info (go 1.18+)
$GO_BUILD_FLAGS = "-trimpath", "-ldflags=-s", "-ldflags=-w", "-buildvcs=false"

# Clean and create bin directory
if (Test-Path $BUILD_DIR) {
    Remove-Item -Path $BUILD_DIR -Recurse -Force
}
New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null

Push-Location src

# Build for Linux AMD64
Write-Host "Building for Linux AMD64..."
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
& go build @GO_BUILD_FLAGS -o "..\$BUILD_DIR\$BINARY_NAME-linux-amd64" .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

# Build for macOS ARM64
Write-Host "Building for macOS ARM64..."
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
& go build @GO_BUILD_FLAGS -o "..\$BUILD_DIR\$BINARY_NAME-darwin-arm64" .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

# Build for Windows AMD64
Write-Host "Building for Windows AMD64..."
$env:GOOS = "windows"
$env:GOARCH = "amd64"
& go build @GO_BUILD_FLAGS -o "..\$BUILD_DIR\$BINARY_NAME-windows-amd64.exe" .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Pop-Location

Write-Host "✓ Binaries built in $BUILD_DIR/"
