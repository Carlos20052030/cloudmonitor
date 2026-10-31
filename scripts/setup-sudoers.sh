#!/bin/bash

# Setup sudoers rule for carlos to run docker without password
# This script should be run with sudo on the VPS

set -euo pipefail

echo "Setting up sudoers rule for docker..."

echo "carlos ALL=(ALL) NOPASSWD: /usr/bin/docker, /usr/local/bin/docker, /usr/bin/docker-compose, /usr/local/bin/docker-compose" | sudo tee /etc/sudoers.d/carlos-docker > /dev/null
sudo chmod 0440 /etc/sudoers.d/carlos-docker

echo "Sudoers rule created successfully"
cat /etc/sudoers.d/carlos-docker