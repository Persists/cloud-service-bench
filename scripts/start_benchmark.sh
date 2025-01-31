#!/bin/sh

N_GENERATORS=$1
N_SINKS=$2

if [ -z "$N_GENERATORS" ] || [ -z "$N_SINKS" ]; then
    echo "Usage: $0 <n_generators> <n_sinks> [zone]"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/wait_startup.sh"

START_TIME=$(date -u -v+2M +"%Y-%m-%dT%H:%M:%SZ")

echo "Starting benchmark at $START_TIME"

for i in $(seq -f "%02g" 1 $N_SINKS); do
    echo "Starting sink-$i"
    gcloud compute scp ./config/experiment/config.yml sink-$i:~/config.yml --zone=europe-west3-c
    gcloud compute ssh sink-$i --zone=europe-west3-c --command "sudo mv ~/config.yml /csb/cloud-service-bench/config/experiment/config.yml"
    gcloud compute ssh sink-$i --zone=europe-west3-c --command "cd /csb/cloud-service-bench; sudo bash -c './sink --instance-name=sink-$i --start-at=\"${START_TIME}\" &> /var/log/benchmark.log &'"
    echo "Started sink-$i"
done

echo "Starting fluentd-sut"
gcloud compute scp ./config/experiment/config.yml fluentd-sut:~/config.yml --zone=europe-west3-c
gcloud compute ssh fluentd-sut --zone=europe-west3-c --command "sudo mv ~/config.yml /csb/cloud-service-bench/config/experiment/config.yml"
gcloud compute ssh fluentd-sut --zone=europe-west3-c --command \
    "cd /csb/cloud-service-bench; sudo bash -c './monitor --instance-name=fluentd-sut --start-at=\"${START_TIME}\" &> /var/log/benchmark.log &'"
echo "Started fluentd-sut"


for i in $(seq -f "%02g" 1 $N_GENERATORS); do
    echo "Starting generator-$i"
    gcloud compute scp ./config/experiment/config.yml generator-$i:~/config.yml --zone=europe-west3-c
    gcloud compute ssh generator-$i --zone=europe-west3-c --command "sudo mv ~/config.yml /csb/cloud-service-bench/config/experiment/config.yml"
    gcloud compute ssh generator-$i --zone=europe-west3-c --command "cd /csb/cloud-service-bench; sudo bash -c './generator --instance-name=generator-$i --start-at=\"${START_TIME}\" &> /var/log/benchmark.log &'"
    echo "Started generator-$i"
done

echo "Benchmark should start at $START_TIME"

