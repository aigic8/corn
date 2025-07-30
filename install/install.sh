#!/usr/bin/env bash

set -e

# TODO: check if required tools exist: tar, curl, ...

# Config 
OWNER="aigic8"
REPO="corn"
BINARY_NAME="corn"
INSTALL_DIR="$HOME/.local/bin"
TMP_DIR=$(mktemp -d)
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)


# Normalize architecture names
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
esac

# Check if OS is linux
if [ "$OS" != "linux" ]; then
	echo "❌ ERROR: Installation script is only functional for linux for now. (identified os: $OS)"
	exit 1
fi

# Fetch Latest Release Info
echo "- Fetching latest release..."
API_URL="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
ASSET_URL=$(curl -s "$API_URL" \
  | grep "browser_download_url" \
  | grep "$OS" \
  | grep "$ARCH" \
  | grep -Eo 'https://[^"]+')

if [ -z "$ASSET_URL" ]; then
  echo "❌ ERROR: No compatible release asset found for $OS/$ARCH."
  exit 1
fi

echo "- Downloading from: $ASSET_URL"
cd "$TMP_DIR"
curl -sLO "$ASSET_URL"

# Extract or Install 
ARCHIVE_FILE=$(basename "$ASSET_URL")

tar -xzf "$ARCHIVE_FILE"

if [ ! -f "$BINARY_NAME" ]; then
  # Try to find it in extracted folder
  BINARY_PATH=$(find . -type f -name "$BINARY_NAME" -perm +111 | head -n 1)
  if [ -z "$BINARY_PATH" ]; then
    echo "❌ ERROR: Binary not found in archive."
    exit 1
  fi
  mv "$BINARY_PATH" "$BINARY_NAME"
fi

chmod +x "$BINARY_NAME"
mv "$BINARY_NAME" "$INSTALL_DIR/"

echo "✅ Installed $BINARY_NAME to $INSTALL_DIR"
echo -e "\t- Make sure installation dir is in your shell path (edit your .bashrc, .zshrc, etc...)"
