sudo apt update
sudo apt-get upgrade -y
curl -fsSL https://toolbelt.treasuredata.com/sh/install-ubuntu-focal-fluent-package5.sh | sh
sudo systemctl start fluentd.service
