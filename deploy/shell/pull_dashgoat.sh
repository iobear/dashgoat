#!/bin/bash

# Define variables
REPO="iobear/dashgoat"
API_URL="https://api.github.com/repos/$REPO/releases/latest"
DEFAULT_BINARY_PATH="./dashgoat"
DOWNLOAD_URL="https://github.com/iobear/dashgoat/releases/download"

# Check if a command line argument is provided for the destination path
if [ -z "$1" ]; then
  BINARY_PATH=$DEFAULT_BINARY_PATH
else
  BINARY_PATH="$1"
fi

# Fetch the latest release information from GitHub API
response=$(curl -s $API_URL)

# Extract the tag name (version) using jq
latest_release=$(echo $response | jq -r '.tag_name' | tr -d '\n' | sed 's/^ *//;s/ *$//')

# Get version information from the binary using strings and grep, and take only the first match
binary_version=$(strings $BINARY_PATH | grep -oP "main.Version=\K[^\']+" | head -n 1 | tr -d '\n' | sed 's/^ *//;s/ *$//')

# Check if binary_version is empty, meaning the binary might not have version info embedded
if [ -z "$binary_version" ]; then
  echo "Could not determine the version of the existing binary."
  binary_version="none"
fi

echo "Your binary version: $binary_version"
echo "Latest release version: $latest_release"

# Compare versions
if [ "$latest_release" == "$binary_version" ]; then
  echo "Your binary is up to date with version $binary_version."
else
  echo "Your binary version $binary_version is not up to date. Latest version is $latest_release."
  echo "Downloading the latest version..."

  # Construct the download URL for the latest binary
  latest_binary_url="$DOWNLOAD_URL/$latest_release/dashgoat"

  # Download the latest binary
  curl -L -o $BINARY_PATH $latest_binary_url

  # Verify the download
  if [ $? -eq 0 ]; then
    chmod +x $BINARY_PATH
    echo "Successfully downloaded and updated to the latest version: $latest_release."
  else
    echo "Failed to download the latest version. Please check the URL and try again."
  fi
fi
