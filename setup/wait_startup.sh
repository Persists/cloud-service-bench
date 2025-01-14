#!/bin/sh

wait_startup_script_to_finish() {
    vm_name=$1
    vm_zone=$2
    echo -n "Waiting for startup script to finish on $vm_name"
    s=""
    while [[ -z "$s" ]]
    do
        sleep 3
        echo -n "."
        s=$(gcloud compute ssh "$vm_name" --zone="$vm_zone" --ssh-flag="-q" --command 'grep -m 1 "startup-script exit status" /var/log/syslog' 2>&1)
        if echo "$s" | grep -q "ERROR: (gcloud.compute.ssh) Could not fetch resource"; then
            echo ""
            echo "Error: VM instance '$vm_name' not found in zone '$vm_zone'."
            return 1
        fi
    done
    echo ""
    return 0
}
