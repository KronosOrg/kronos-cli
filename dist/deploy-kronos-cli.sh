#!/bin/bash

# Build the Go app
echo "Downloading the binaries..."
curl -LO https://storage.cloud.google.com/kronos-cli/$(curl -L -s https://storage.googleapis.com/kronos-cli/stable.txt)/kronos-cli

# Add executable permissions
chmod +x kronos-cli

# Move the executable to /usr/local/bin
sudo mv kronos-cli /usr/local/bin

echo "Kronos CLI is Downloaded and Installed in the OS"
