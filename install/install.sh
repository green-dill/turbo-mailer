#!/bin/bash
set -e

if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
elif [ -f /etc/lsb-release ]; then
    . /etc/lsb-release
    OS=$DISTRIB_ID
else
    echo "Unsupported operating system"
    exit 1
fi

# Check if the script is run as root
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# Update and install packages based on the detected OS
case "$OS" in
    "Ubuntu"|"Debian GNU/Linux")
        apt-get update
        apt-get install -y curl git unzip
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# ipv6 issue: https://github.com/oven-sh/bun/issues/10698
# Disable IPv6
sysctl -w net.ipv6.conf.all.disable_ipv6=1
sysctl -w net.ipv6.conf.default.disable_ipv6=1
sysctl -w net.ipv6.conf.lo.disable_ipv6=1

# Make IPv6 settings persistent
echo "net.ipv6.conf.all.disable_ipv6 = 1" >> /etc/sysctl.conf
echo "net.ipv6.conf.default.disable_ipv6 = 1" >> /etc/sysctl.conf
echo "net.ipv6.conf.lo.disable_ipv6 = 1" >> /etc/sysctl.conf

# Apply sysctl changes
sysctl -p


# Install bun
curl -fsSL https://bun.sh/install | bash

# Add bun to PATH for this session
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# Clone the repository
# This Personal Access Token (PAT) will expire on 2025-10-01 (GitHub requires PATs to have an expiration date within 1 year)
# Check if the PAT is expired or about to expire
PAT_EXPIRATION="2025-10-01"
CURRENT_DATE=$(date +%Y-%m-%d)
DAYS_UNTIL_EXPIRATION=$(( ($(date -d "$PAT_EXPIRATION" +%s) - $(date -d "$CURRENT_DATE" +%s)) / 86400 ))
if [ $DAYS_UNTIL_EXPIRATION -le 0 ]; then
    echo "ERROR: The GitHub Personal Access Token has expired. Please update the token in the script."
    exit 1
elif [ $DAYS_UNTIL_EXPIRATION -le 60 ]; then
    echo "WARNING: The GitHub Personal Access Token will expire in $DAYS_UNTIL_EXPIRATION days. Please update the token soon."
fi

git clone https://yinheli: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@github.com/green-dill/turbo-mailer.git --depth 1

# Navigate to the install directory
cd turbo-mailer/install

# Install dependencies
bun install

# Run the installation script
bun index.ts
