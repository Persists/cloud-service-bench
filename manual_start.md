This is a manual start guide for fluentd. This is just a first step to get started.
TODO: Add a more automated way to start fluentd.

```bash
    sudo apt update
    sudo apt-get upgrade -y

    # Install fluent-package5
    curl -fsSL https://toolbelt.treasuredata.com/sh/install-ubuntu-noble-fluent-package5.sh | sh

    # Start fluentd
    sudo systemctl start fluentd.service

    # Check fluentd status
    sudo systemctl status fluentd.service

    # Stop fluentd
    sudo systemctl stop fluentd.service
```
