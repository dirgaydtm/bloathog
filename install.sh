#!/bin/sh
# Bloathog - Installation Script
# Downloads the latest release binary for Linux/macOS

set -e

echo "Installing Bloathog..."

# Detect OS
OS="$(uname -s)"

# Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="x86_64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  i386|i686)     ARCH="i386" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Latest release URL
RELEASE_URL="https://github.com/dirgaa/bloathog/releases/latest/download"
TAR_FILE="bloathog_${OS}_${ARCH}.tar.gz"

echo "Downloading ${TAR_FILE}..."

# Download and Extract
curl -sL "${RELEASE_URL}/${TAR_FILE}" -o "/tmp/${TAR_FILE}"
tar -xzf "/tmp/${TAR_FILE}" -C /tmp bloathog
chmod +x /tmp/bloathog

# Install to PATH
if command -v sudo >/dev/null 2>&1; then
  sudo install /tmp/bloathog /usr/local/bin/bloathog
  echo "✅ Installed to /usr/local/bin/bloathog"
else
  install /tmp/bloathog /usr/local/bin/bloathog
  echo "✅ Installed to /usr/local/bin/bloathog"
fi

# Cleanup
rm -f "/tmp/${TAR_FILE}" /tmp/bloathog

echo ""
echo "Bloathog installed successfully!"
echo "Usage: bloathog"
echo "       bloathog <your command>"
