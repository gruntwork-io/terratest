#!/usr/bin/env bash
set -euo pipefail

# Send the log output from this script to user-data.log, syslog, and the console
# From: https://alestic.com/2010/12/ec2-user-data-output/
exec > >(tee /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1

# Create our new 'terratest' user
adduser --disabled-password --gecos "" terratest

# Create CA pubkey file
mkdir -p /etc/ssh
cat > /etc/ssh/trusted-user-ca-keys.pub <<'EOKEY'
${ssh_ca_public_key}
EOKEY

# Drop-in configuration for sshd
mkdir -p /etc/ssh/sshd_config.d
echo 'TrustedUserCAKeys /etc/ssh/trusted-user-ca-keys.pub' > /etc/ssh/sshd_config.d/ca.conf

# Bounce the service to apply the config change
service ssh restart
