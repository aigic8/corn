#!/usr/bin/env sh

set -e

# TODO: check if required tools exist: tar, curl, ...

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[1;34m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m' # No Color

PREF="$BLUE■$NC"
ERR_PREF="❌$RED ERROR:$NC"
SUCCESS_PREF="$GREEN■$NC"

# Config 
OWNER="aigic8"
REPO="corn"
BINARY_NAME="corn"
INSTALL_DIR="$HOME/.local/bin"
SERVICE_PATH="/etc/systemd/system/corn.service"
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
	echo -e "$ERR_PREF Installation script is only functional for linux for now. (identified os: $OS)"
	exit 1
fi

echo "Welcome to corn installer!"

# empty new line
echo 

# Fetch Latest Release Info
echo -e "$PREF Fetching latest release..."
API_URL="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
ASSET_URL=$(curl -s "$API_URL" \
  | grep "browser_download_url" \
  | grep "$OS" \
  | grep "$ARCH" \
  | grep -Eo 'https://[^"]+')

if [ -z "$ASSET_URL" ]; then
  echo -e "$ERR_PREF No compatible release asset found for $OS/$ARCH."
  exit 1
fi

echo -e "$PREF Downloading from: $BOLD$ASSET_URL$NC"
cd "$TMP_DIR"
curl -sLO "$ASSET_URL"

# Extract or Install 
ARCHIVE_FILE=$(basename "$ASSET_URL")

tar -xzf "$ARCHIVE_FILE"

if [ ! -f "$BINARY_NAME" ]; then
  # Try to find it in extracted folder
  BINARY_PATH=$(find . -type f -name "$BINARY_NAME" -perm +111 | head -n 1)
  if [ -z "$BINARY_PATH" ]; then
    echo -e "$ERR_PREF Binary not found in archive."
    exit 1
  fi
  mv "$BINARY_PATH" "$BINARY_NAME"
fi

chmod +x "$BINARY_NAME"
mv "$BINARY_NAME" "$INSTALL_DIR/"

# TODO: add option to choose the installation path

echo -e "\t$SUCCESS_PREF Installed $BINARY_NAME to $BOLD$INSTALL_DIR$NC"

# empty new line
echo 

# Add Systemd Service
# TODO: add option to choose the user (now it is the current user)
echo -e "$PREF Adding Systemd service to $BLUE$SERVICE_PATH$NC $YELLOW(requires root access)$NC"
curl -s https://raw.githubusercontent.com/aigic8/corn/refs/heads/main/install/corn.service.template \
  | sed -e "s|\$\$INSTALL_BIN|$INSTALL_DIR/$BINARY_NAME|" \
        -e "s|\$\$USERNAME|$USER|" \
	| sudo tee $SERVICE_PATH > /dev/null

echo -e "\t$SUCCESS_PREF Added Systemd service $BOLD(corn.service)$NC"

echo 
echo -e "$PREF To finalize the installation do the following:"
echo -e "\t$SUCCESS_PREF Make sure installation dir is in your shell path (edit your .bashrc, .zshrc, etc...) $BLUE(REQUIRED)$NC"
echo -e "\t$SUCCESS_PREF Start the service using: sudo systemctl start corn.service $BLUE(REQUIRED)$NC"
echo -e "\t$SUCCESS_PREF Enable the service using (the service will start on reboot): sudo systemctl enable corn.service $YELLOW(OPTIONAL)$NC"
echo -e "\t$SUCCESS_PREF View the status of the service: sudo systemctl status corn.service"
echo -e "\t$SUCCESS_PREF View the logs of the service: sudo journalctl -u corn.service"

# empty new line
echo
