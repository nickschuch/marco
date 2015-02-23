#!/bin/bash

# Name:        provision.sh
# Author:      Nick Schuch
# Description: Provisions an environment for testing Golang.

# Install Golang so we can compile and run.
GOBALL="go1.4.1.linux-amd64.tar.gz"

echo "Download and build ${GOBALL}"
cd /tmp
wget -q https://storage.googleapis.com/golang/${GOBALL}
tar -C /usr/local -xzf ${GOBALL}
chown -R vagrant:vagrant /usr/local/go
mkdir -p /opt/golang
chmod -R 777 /opt/golang

echo 'GOPATH=/opt/golang' >> /etc/environment
echo 'PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo 'PATH=$PATH:/opt/golang/bin' >> /etc/profile
export GOPATH=/opt/golang
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/opt/golang/bin

# Some other random packages.
export DEBIAN_FRONTEND=noninteractive
sudo -E apt-get install -y vim make git
