#!/bin/sh

N_GENERATORS=$1
N_SINKS=$2

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/wait_startup.sh"

if [ -z "$N_GENERATORS" ] || [ -z "$N_SINKS" ]; then
    echo "Usage: $0 <n_generators> <n_sinks>"
    exit 1
fi

for i in $(seq -f "%02g" 1 $N_SINKS); do
    if ! wait_startup_script_to_finish "sink-$i" "europe-west3-c"; then
        echo "Startup script failed or VM not found. Exiting."
        exit 1
    fi
done

# Check all instances to finish startup script
for i in $(seq -f "%02g" 1 $N_GENERATORS); do
    if ! wait_startup_script_to_finish "generator-$i" "europe-west3-c"; then
        echo "Startup script failed or VM not found. Exiting."
        exit 1
    fi
done