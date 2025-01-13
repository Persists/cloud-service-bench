#!/bin/sh

gcloud compute scp ./config/fluentd/fluentd.conf victoria-r1:~/fluentd.conf --zone=europe-west3-c
gcloud compute ssh victoria-r1 --zone=europe-west3-c --command "sudo mv fluentd.conf /etc/fluent/ && sudo systemctl restart fluentd"

