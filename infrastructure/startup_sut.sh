sudo apt update -y
sudo apt install -y curl wget git htop

wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
echo 'Defaults        secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin:/usr/local/go/bin"' | sudo tee /etc/sudoers.d/spath

curl -fsSL https://toolbelt.treasuredata.com/sh/install-ubuntu-focal-fluent-package5.sh | sh
sudo systemctl start fluentd.service

mkdir csb
cd csb

git clone https://github.com/Persists/cloud-service-bench.git
cd cloud-service-bench
sudo go build  -o ./monitor ./cmd/monitor/main.go

echo done