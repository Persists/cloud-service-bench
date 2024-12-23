sudo apt update -y
sudo apt upgrade -y
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >>~/.profile

# TODO: compile the go program instead of using go run