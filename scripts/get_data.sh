#!/bin/sh

N_GENERATORS=$1
N_SINKS=$2


if [ -z "$N_GENERATORS" ] || [ -z "$N_SINKS" ]; then
    echo "Usage: $0 <n_generators> <n_sinks> [zone]"
    exit 1
fi

# Logs
mkdir -p "./results/logs"
for i in $(seq -f "%02g" 1 $N_SINKS); do
    gcloud compute scp sink-$i:/var/log/benchmark.log ./results/logs/sink-$i.log --zone=europe-west3-c
done

gcloud compute scp fluentd-sut:/var/log/benchmark.log ./results/logs/fluentd-monitor.log --zone=europe-west3-c
gcloud compute scp fluentd-sut:/var/log/fluent/fluentd.log ./results/logs/fluentd-fluentd.log --zone=europe-west3-c

for i in $(seq -f "%02g" 1 $N_GENERATORS); do
    gcloud compute scp generator-$i:/var/log/benchmark.log ./results/logs/generator-$i.log --zone=europe-west3-c
done

mkdir -p "./results/data"
gcloud compute scp --recurse fluentd-sut:/csb/cloud-service-bench/results ./results/data/fluentd-sut --zone=europe-west3-c
 

# Data
for i in $(seq -f "%02g" 1 $N_SINKS); do
    gcloud compute scp  --recurse sink-$i:/csb/cloud-service-bench/results ./results/data/sink-$i --zone=europe-west3-c
done

terraform -chdir=infrastructure destroy -auto-approve

terraform -chdir=infrastructure destroy -auto-approve -target="google_compute_instance_from_template.generator[0]" -target="google_compute_instance_from_template.generator[1]" -target="google_compute_instance_from_template.generator[2]" -target="google_compute_instance_from_template.generator[3]"
terraform -chdir=infrastructure destroy -auto-approve -target="google_compute_instance.fluentd" 









