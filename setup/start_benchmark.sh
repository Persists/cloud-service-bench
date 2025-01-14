#!/bin/sh

N_GENERATORS=$1
N_SINKS=$2

if [ -z "$N_GENERATORS" ] || [ -z "$N_SINKS" ]; then
    echo "Usage: $0 <n_generators> <n_sinks> [zone]"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/wait_startup.sh"

for i in $(seq -f "%02g" 1 $N_SINKS); do
    gcloud compute scp ./config/experiment/config.yml sink-$i:~/config.yml --zone=europe-west3-c
    gcloud compute ssh sink-$i --zone=europe-west3-c --command "sudo mv ~/config.yml /csb/cloud-service-bench/config/experiment/config.yml"
    gcloud compute ssh sink-$i --zone=europe-west3-c --command "cd /csb/cloud-service-bench; sudo bash -c './sink --instance-name=sink-$i &> /var/log/benchmark.log &'"
done

gcloud compute ssh fluentd-sut --zone=europe-west3-c --command \
    "cd /csb/cloud-service-bench; sudo bash -c './monitor --instance-name=fluentd-sut &> /var/log/benchmark.log &'"

START_TIME=$(date -u -v+0M +"%Y-%m-%dT%H:%M:%SZ")

for i in $(seq -f "%02g" 1 $N_GENERATORS); do
    gcloud compute scp ./config/experiment/config.yml generator-$i:~/config.yml --zone=europe-west3-c
    gcloud compute ssh generator-$i --zone=europe-west3-c --command "sudo mv ~/config.yml /csb/cloud-service-bench/config/experiment/config.yml"
    gcloud compute ssh generator-$i --zone=europe-west3-c --command "cd /csb/cloud-service-bench; sudo bash -c './generator --instance-name=generator-$i --start-at=\"${START_TIME}\" &> /var/log/benchmark.log &'"
done
done

echo "Start scripts finished successfully."
