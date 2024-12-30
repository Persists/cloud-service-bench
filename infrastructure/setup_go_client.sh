#! /bin/bash
apt update -y
apt upgrade -y
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >>~/.profile
source ~/.profile