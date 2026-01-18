#!/usr/bin/env bash
set -euo pipefail

REPO="subcode-labs/dots"
BIN_NAME="dots"
INSTALL_DIR="/usr/local/bin"
FALLBACK_DIR="$HOME/bin"

color() {
  local code="$1"
  shift
  printf "\033[%sm%s\033[0m\n" "$code" "$*"
}

info() {
  color "1;34" "$*"
}

success() {
  color "1;32" "$*"
}

warn() {
  color "1;33" "$*"
}

error() {
  color "1;31" "$*"
  exit 1
}

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux|darwin) ;;
  *) error "Unsupported OS: $OS" ;;
 esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) error "Unsupported architecture: $ARCH" ;;
 esac

info "Detected $OS/$ARCH"

if ! command -v curl >/dev/null 2>&1; then
  error "curl is required but not installed"
fi

LATEST_TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep -m 1 '"tag_name"' | cut -d '"' -f 4)"
if [[ -z "$LATEST_TAG" ]]; then
  error "Unable to find latest release"
fi

VERSION="${LATEST_TAG#v}"
ARCHIVE_NAME="${BIN_NAME}_${VERSION}_${OS}_${ARCH}"
ARCHIVE_FILE="${ARCHIVE_NAME}.tar.gz"
if [[ "$OS" == "windows" ]]; then
  ARCHIVE_FILE="${ARCHIVE_NAME}.zip"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ARCHIVE_FILE}"

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

info "Downloading ${DOWNLOAD_URL}"
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE_FILE"

info "Extracting archive"
if [[ "$ARCHIVE_FILE" == *.zip ]]; then
  if ! command -v unzip >/dev/null 2>&1; then
    error "unzip is required to extract zip archives"
  fi
  unzip -q "$TMP_DIR/$ARCHIVE_FILE" -d "$TMP_DIR"
else
  tar -xzf "$TMP_DIR/$ARCHIVE_FILE" -C "$TMP_DIR"
fi

TARGET_DIR="$INSTALL_DIR"
if [[ ! -w "$INSTALL_DIR" ]]; then
  if command -v sudo >/dev/null 2>&1; then
    info "Installing to $INSTALL_DIR (sudo required)"
    sudo install -m 0755 "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
    success "Installed $BIN_NAME to $INSTALL_DIR"
    exit 0
  fi
  warn "No write access to $INSTALL_DIR, falling back to $FALLBACK_DIR"
  mkdir -p "$FALLBACK_DIR"
  TARGET_DIR="$FALLBACK_DIR"
fi

install -m 0755 "$TMP_DIR/$BIN_NAME" "$TARGET_DIR/$BIN_NAME"

success "Installed $BIN_NAME to $TARGET_DIR"
if [[ "$TARGET_DIR" == "$FALLBACK_DIR" ]]; then
  warn "Ensure $FALLBACK_DIR is on your PATH"
fi
