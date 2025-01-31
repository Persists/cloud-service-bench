#!/bin/sh

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/wait_startup.sh"

if wait_startup_script_to_finish "fluentd-sut" "europe-west3-c"; then
    echo "Startup script finished successfully."

    gcloud compute scp ./config/fluentd/fluentd.conf fluentd-sut:~/fluentd.conf --zone=europe-west3-c
    gcloud compute ssh fluentd-sut --zone=europe-west3-c --command "sudo mv fluentd.conf /etc/fluent/ && sudo systemctl restart fluentd"
else
    echo "Startup script failed or VM not found. Exiting."
    exit 1
fi
