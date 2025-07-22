#!/usr/bin/env bash

# Installs the "duf" and "muffet" admin tools on Linux.
# bash <(curl -s https://raw.githubusercontent.com/photoprism/photoprism/develop/scripts/dist/install-admin-tools.sh)

# Abort if not executed as root..
if [[ $(id -u) != "0" ]]; then
  echo "Usage: run ${0##*/} as root" 1>&2
  exit 1
fi

set -eux;

# Is Go installed?
if ! command -v go &> /dev/null
then
    echo "Go must be installed to run this."
    exit 1
fi

echo "Installing the duf command to check storage usage..."
GOBIN="/usr/local/bin" go install github.com/muesli/duf@latest

echo "Installing muffet, a tool for checking links..."
GOBIN="/usr/local/bin" go install github.com/raviqqe/muffet@latest

echo "Installing petname to generate pronounceable names..."
GOBIN="/usr/local/bin" go install github.com/dustinkirkland/golang-petname/cmd/petname@latest

echo "Installing doctl for using the DigitalOcean API...."
GOBIN="/usr/local/bin" go install github.com/digitalocean/doctl/cmd/doctl@latest

echo "Installing Kustomize for Kubernetes configuration management...."
GOBIN="/usr/local/bin" go install sigs.k8s.io/kustomize/kustomize/v5@latest

# Create a symbolic link for "duf" so that it is used instead of the original "df".
ln -sf /usr/local/bin/duf /usr/local/bin/df