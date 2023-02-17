#!/bin/bash
# This script is meant to be run in the User Data of an EC2 Instance while it's booting. It starts a Ruby web app.
# This script assumes it is running in an AMI built from the Packer templates in examples/packer-docker-example
# (either build.json or build.pkr.hcl).

set -e

# Send the log output from this script to user-data.log, syslog, and the console
# From: https://alestic.com/2010/12/ec2-user-data-output/
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

# The variables below are filled in using Terraform interpolation
nohup ruby /home/ubuntu/app.rb "${instance_port}" "${instance_text}" &
