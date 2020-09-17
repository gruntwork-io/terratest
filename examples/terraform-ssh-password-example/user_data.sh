#!/usr/bin/env bash
set -euo pipefail

# Send the log output from this script to user-data.log, syslog, and the console
# From: https://alestic.com/2010/12/ec2-user-data-output/
exec > >(tee /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1

# Create our new 'terratest' user
adduser --disabled-password --gecos "" terratest

# Set the user's password based on the random input from 'test/terraform_ssh_password_example_test.go'
# shellcheck disable=SC2154
echo "terratest:${terratest_password}" | chpasswd

# Enable password auth on the SSH service
sed -i 's/^PasswordAuthentication no$/PasswordAuthentication yes/g' /etc/ssh/sshd_config

# Bounce the service to apply the config change
service ssh restart
