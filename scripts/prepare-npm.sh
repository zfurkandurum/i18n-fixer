#!/bin/bash
set -euo pipefail

# This script prepares npm platform packages from GoReleaser output.
# Called during the release CI pipeline after GoReleaser creates binaries.

VERSION="${1:?Usage: prepare-npm.sh <version>}"
DIST_DIR="${2:-dist}"

PLATFORMS=(
  "darwin-arm64:darwin_arm64"
  "darwin-x64:darwin_amd64_v1"
  "linux-x64:linux_amd64_v1"
  "linux-arm64:linux_arm64"
  "win32-x64:windows_amd64_v1"
  "win32-arm64:windows_arm64"
)

for entry in "${PLATFORMS[@]}"; do
  NPM_PLATFORM="${entry%%:*}"
  GO_PLATFORM="${entry##*:}"
  PKG_NAME="i18n-fixer-${NPM_PLATFORM}"
  PKG_DIR="npm/${PKG_NAME}"

  echo "Preparing ${PKG_NAME}..."

  mkdir -p "${PKG_DIR}/bin"

  # Determine binary name
  BIN_NAME="i18n-fixer"
  if [[ "${NPM_PLATFORM}" == win32-* ]]; then
    BIN_NAME="i18n-fixer.exe"
  fi

  # Copy binary from GoReleaser output
  SRC="${DIST_DIR}/i18n-fixer_${GO_PLATFORM}/${BIN_NAME}"
  if [ ! -f "${SRC}" ]; then
    echo "  Warning: ${SRC} not found, skipping"
    continue
  fi

  cp "${SRC}" "${PKG_DIR}/bin/${BIN_NAME}"
  chmod +x "${PKG_DIR}/bin/${BIN_NAME}"

  # Determine OS and CPU for package.json
  OS="${NPM_PLATFORM%%-*}"
  CPU="${NPM_PLATFORM##*-}"
  if [ "${CPU}" = "x64" ]; then CPU_JSON="x64"; fi
  if [ "${CPU}" = "arm64" ]; then CPU_JSON="arm64"; fi

  # Create package.json
  cat > "${PKG_DIR}/package.json" << PKGJSON
{
  "name": "${PKG_NAME}",
  "version": "${VERSION}",
  "description": "i18n-fixer binary for ${NPM_PLATFORM}",
  "os": ["${OS}"],
  "cpu": ["${CPU}"],
  "bin": {
    "i18n-fixer": "bin/${BIN_NAME}"
  },
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/i18n-fixer/i18n-fixer.git"
  }
}
PKGJSON

  echo "  Done: ${PKG_DIR}"
done

# Update version in main package
sed -i.bak "s/\"version\": \".*\"/\"version\": \"${VERSION}\"/" npm/i18n-fixer/package.json
rm -f npm/i18n-fixer/package.json.bak

# Update optional dependency versions
for entry in "${PLATFORMS[@]}"; do
  NPM_PLATFORM="${entry%%:*}"
  PKG_NAME="i18n-fixer-${NPM_PLATFORM}"
  sed -i.bak "s/\"${PKG_NAME}\": \".*\"/\"${PKG_NAME}\": \"${VERSION}\"/" npm/i18n-fixer/package.json
  rm -f npm/i18n-fixer/package.json.bak
done

echo "All npm packages prepared for version ${VERSION}"
