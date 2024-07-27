#!/bin/bash

REPO_OWNER="WildEgor"
REPO_NAME="pi-storyteller"
FILE_PATTERN="linux_armv6.tar.gz" # linux_amd.tar.gz
APP_DIR="/app"
APP_BIN_PATH="/app/bin/pi-storyteller"
APP_CONFIG_PATH="/app/config.yaml"
UNIT_USER="wild"
UNIT_NAME="pi_storyteller.service"
UNIT_FILE_PATH="/etc/systemd/system/pi_storyteller.service"
TEMP_DIR="/app/temp"

cleanup() {
    echo "Cleaning up..."
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

echo "Install deps..."
apt-get update && apt-get install -y bash curl libxml2-utils jq

echo "Check/create unit for service..."
if [ -f "$UNIT_FILE_PATH" ]; then
    echo "Unit file already exists: $UNIT_FILE_PATH"
else
    echo "Creating unit file: $UNIT_FILE_PATH"
    cat <<EOF | sudo tee "$UNIT_FILE_PATH"
[Unit]
Description=Pi Storyteller Service

[Install]
WantedBy=multi-user.target

[Service]
ExecStart=$APP_BIN_PATH
WorkingDirectory=$APP_DIR
User=$UNIT_USER
Restart=always
RestartSec=5
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=pi_storyteller
EOF
    systemctl daemon-reload
    echo "Unit file created and systemd reloaded."
fi

# Fetch latest release
echo "Fetching latest release..."
RELEASE_URL=$(curl -s https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest | jq -r --arg pattern "$FILE_PATTERN" '.assets[] | select(.name | test($pattern)) | .browser_download_url')
if [ -z "$RELEASE_URL" ]; then
    echo "Failed to extract the release URL. Exiting."
    exit 1
fi
mkdir -p "$TEMP_DIR"
curl -L "$RELEASE_URL" -o "$TEMP_DIR/latest_release.tar.gz"
tar -xzf "$TEMP_DIR/latest_release.tar.gz" -C "$TEMP_DIR" --wildcards 'bin'
mkdir -p "$APP_DIR"
if [ -f "$TEMP_DIR/bin" ]; then
    mv "$TEMP_DIR/bin" "$APP_DIR"
else
    echo "$TEMP_DIR/bin directory not found in the extracted contents."
    exit 1
fi

echo "Prepare assets..."
SOURCE_CODE_URL=$(curl -sL https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/tags | jq -r '.[0].tarball_url')
curl -sL "$SOURCE_CODE_URL" -o "$TEMP_DIR/latest_source.tar.gz"
tar -xzf "$TEMP_DIR/latest_source.tar.gz" -C "$TEMP_DIR"
EXTRACTED_DIR=$(tar -tf "$TEMP_DIR/latest_source.tar.gz" | head -1 | cut -f1 -d'/')
mkdir -p "$APP_DIR/assets" && cp -r "$TEMP_DIR/$EXTRACTED_DIR/assets" "$APP_DIR/assets"
mkdir -p "$APP_DIR/scripts" && cp -r "$TEMP_DIR/$EXTRACTED_DIR/scripts" "$APP_DIR/scripts"

echo "Restart app..."
export PI_STORYTELLER_CONFIG_PATH=$APP_CONFIG_PATH
systemctl restart $UNIT_NAME
# systemctl status $UNIT_NAME
sudo journalctl UNIT=$UNIT_NAME
exit 1
