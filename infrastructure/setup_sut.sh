### Install Go
sudo apt update -y
sudo apt upgrade -y
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile

curl -fsSL https://toolbelt.treasuredata.com/sh/install-ubuntu-focal-fluent-package5.sh | sh
sudo systemctl start fluentd.service

