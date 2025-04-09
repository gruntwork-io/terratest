#!/bin/bash
# Install and configure a simple web app built on top of Ruby and Sinatra

set -e

readonly APP_RB_SRC="/tmp/packer-docker-example/app.rb"
readonly APP_RB_DST="/home/ubuntu/app.rb"

echo "Installing Ruby"
sudo apt-get update
sudo apt-get install -y make zlib1g-dev build-essential ruby ruby-dev

echo "Installing Sinatra"
sudo gem install sinatra json rackup puma

echo "Moving $APP_RB_SRC to $APP_RB_DST"
mkdir -p "$(dirname "$APP_RB_DST")"
mv "$APP_RB_SRC" "$APP_RB_DST"