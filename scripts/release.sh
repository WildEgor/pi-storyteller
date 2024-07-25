#!/bin/bash

# Configuration
REPO_OWNER="WildEgor"
REPO_NAME="pi-storyteller"
FILE_PATTERN="pi-storyteller_.*_linux_armv6.tar.gz"
DESTINATION_DIR="/app"

# Fetch latest release info from GitHub API
echo "Fetching latest release info..."
RELEASE_JSON=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest")

# Extract the URL for the desired asset
echo "Searching for the file matching pattern '$FILE_PATTERN'..."
ASSET_URL=$(echo "$RELEASE_JSON" | grep -oP '"browser_download_url": "\K(.*?)(?=")' | grep -E "$FILE_PATTERN")

if [ -z "$ASSET_URL" ]; then
  echo "No file matching pattern '$FILE_PATTERN' found in the latest release."
  exit 1
fi

# Download the specified file
echo "Downloading file from URL: $ASSET_URL"
curl -L "$ASSET_URL" -o latest_release.tar.gz

# Extract the 'bin' directory from the tarball
echo "Extracting 'bin' directory..."
mkdir -p temp_extraction
tar -xzf latest_release.tar.gz -C temp_extraction --wildcards 'bin'

# Move the 'bin' directory to the destination
echo "Moving 'bin' directory to $DESTINATION_DIR..."
if [ -f "temp_extraction/bin" ]; then
  mv temp_extraction/bin "$DESTINATION_DIR"
else
  echo "'bin' directory not found in the extracted contents."
  rm -rf latest_release.tar.gz temp_extraction
  exit 1
fi

# Clean up
echo "Cleaning up..."
rm -rf latest_release.tar.gz temp_extraction

# Restart
echo "Restart app..."
systemctl restart pi_storyteller.service
systemctl status pi_storyteller.service
