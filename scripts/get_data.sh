#!/bin/sh

N_GENERATORS=$1
N_SINKS=$2
EXPERIMENT_ID=$3

if [ -z "$N_GENERATORS" ]; then
    echo "N_GENERATORS is not set"
    exit 1
fi

if [ -z "$N_SINKS" ]; then
    echo "N_SINKS is not set"
    exit 1
fi

if [ -z "$EXPERIMENT_ID" ]; then
    echo "EXPERIMENT_ID is not set"
    exit 1
fi

# Logs
mkdir -p "./results/$EXPERIMENT_ID/logs"
for i in $(seq -f "%02g" 1 $N_SINKS); do
    gcloud compute scp sink-$i:/var/log/benchmark.log ./results/$EXPERIMENT_ID/logs/sink-$i.log --zone=europe-west3-c
done

gcloud compute scp fluentd-sut:/var/log/benchmark.log ./results/$EXPERIMENT_ID/logs/fluentd-monitor.log --zone=europe-west3-c
gcloud compute scp fluentd-sut:/var/log/fluent/fluentd.log ./results/$EXPERIMENT_ID/logs/fluentd-fluentd.log --zone=europe-west3-c

for i in $(seq -f "%02g" 1 $N_GENERATORS); do
    gcloud compute scp generator-$i:/var/log/benchmark.log ./results/$EXPERIMENT_ID/logs/generator-$i.log --zone=europe-west3-c
done

# Data
mkdir -p "./results/$EXPERIMENT_ID/data"
gcloud compute scp --recurse fluentd-sut:/csb/cloud-service-bench/results ./results/$EXPERIMENT_ID/data/fluentd-sut --zone=europe-west3-c

for i in $(seq -f "%02g" 1 $N_SINKS); do
    gcloud compute scp --recurse sink-$i:/csb/cloud-service-bench/results ./results/$EXPERIMENT_ID/data/sink-$i --zone=europe-west3-c
done

terraform -chdir=infrastructure destroy -auto-approve