This is the start guide for fluentd. This is a WIP file, I will extand this file with further information while the benchmark progresses.

### Set up a GCP instance for the SUT

This takes a while because it is also installing fluentd and updating dependencies before

```bash
terraform -chdir=infrastructure apply

# see the progress of the startup script:
gcloud compute ssh victoria-r1 --zone=europe-west3-c
sudo tail -f /var/log/syslog | grep "startup-script"
```

### Reconfigure fluentd

```bash
gcloud compute scp ./config/fluentd/fluentd.conf victoria-r1:~/fluentd.conf --zone=europe-west3-c
gcloud compute ssh victoria-r1 --zone=europe-west3-c --command "sudo mv fluentd.conf /etc/fluent/ && sudo systemctl restart fluentd"
```

### Retrieve the ip

```bash
gcloud compute instances describe victoria-r1 --zone=europe-west3-c --format="get(networkInterfaces[0].accessConfigs[0].natIP)"
```

### Make a request

Therefore first the ip in the golang script needs to be updated. Then the script can be run.

```bash
go run ./cmd/main.go
```

### Inspect remote fluentd logs

```bash
gcloud compute ssh victoria-r1 --zone=europe-west3-c --command "sudo tail -f /var/log/fluent/fluentd.log"
```
