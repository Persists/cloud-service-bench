#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Sat Jan 25 21:00:27 2025

@author: jonasheisterberg
"""

# %%
from collections import defaultdict
from datetime import datetime
import json
import os
from tqdm import tqdm

preprocessed_dir = "./results/preprocessed"

def read_sink_file(file_path):
    '''This function reads a sink log file and extracts the message counts per timestamp.'''
    file_name = os.path.basename(file_path)
    

    def estimate_line_count(file_path):
        with open(file_path, "r") as file:
            return sum(1 for _ in file)
                
    progress_bar = tqdm(total=estimate_line_count(file_path), desc=f"Processing {file_name}", unit="line", dynamic_ncols=True)

    # Read only the relevant lines
    experiment_started = False

    message_counts = defaultdict(int)

    try:
        with open(file_path, "r") as file:
            for line in file:
                progress_bar.update(1)
                
                if "Experiment phase" in line:
                    experiment_started = True
                    continue

                if "Cool-down phase" in line:
                    progress_bar.n = progress_bar.total
                    progress_bar.refresh()
                    break

                # skip empty lines
                if line.strip() == "":
                    continue
                
                if experiment_started:
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError as e:
                        print(f"Error decoding JSON: {e}")
                        continue
                            
                    timestamp_str = data["arrivalTime"]
                    timestamp = datetime.fromisoformat(timestamp_str)
                    message_count = len(data.get("logMessages", []))
                    message_counts[timestamp.isoformat()] += message_count                
    except FileNotFoundError:
        print(f"File not found: {file_name}")
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
    return message_counts

def read_metrics_file(file_path):
    '''This function reads a metrics log file and extracts the metrics per timestamp.'''
    file_name = os.path.basename(file_path)

    metrics = defaultdict(dict)

    experiment_started = False

    try:
        with open(file_path, "r") as file:
            for line in file:
                if "Experiment phase" in line:
                    experiment_started = True
                    continue

                if "Cool-down phase" in line:
                    break

                # skip empty lines
                if line.strip() == "":
                    continue

                if experiment_started:
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError as e:
                        print(f"Error decoding JSON: {e}")
                        continue
                    timestamp_str = data["time"]
                    timestamp = datetime.fromisoformat(timestamp_str.replace("Z", "+00:00"))
                    metrics[timestamp.isoformat()] = data
    except FileNotFoundError:
        print(f"File not found: {file_name}")
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")

    return metrics

def data_from_results_directory(directory_path):
    '''
    This functions accepts a directory path which points to a single experiment directory.
    
    It then reads the data from sink and metrics files and returns a dictionary with the data.
    '''

    sink_directory = os.path.join(directory_path, "data/sink-01")
    metrics_directory = os.path.join(directory_path, "data/fluentd-sut")

    print(sink_directory)
    print(metrics_directory)

    try:
        sink_file = next(
            os.path.join(sink_directory, f) for f in os.listdir(sink_directory) if f.endswith(".log")
        )
        metrics_file = next(
            os.path.join(metrics_directory, f) for f in os.listdir(metrics_directory) if f.endswith(".log")
        )
    except StopIteration:
        raise FileNotFoundError("Required .log files not found in one or both directories.")

    message_counts = read_sink_file(sink_file)
    metrics = read_metrics_file(metrics_file)

    return {
        "file_names": {
            "sink": sink_file,
            "metrics": metrics_file
        },
        "message_counts": message_counts,
        "metrics": metrics
    }

# %%
# This section was used to preprocess the data from the results directory and store it in a preprocessed directory.
# This was done to speed up the later analysis steps. Reading the data from the log files was the most time-consuming
# part of the analysis process.
directory_path = "./results/3.4worker.memory"

worker = 4
experimentId = 4003

data = data_from_results_directory(directory_path)
# add metadata and store in file
storable = {
    "metadata": {
        "worker": worker,
        "expetimentId": experimentId,
        "sink": os.path.basename(data["file_names"]["sink"]),
        "metrics": os.path.basename(data["file_names"]["metrics"]),
    },
    "message_counts": data["message_counts"],
    "metrics": data["metrics"]
}

destination_file = os.path.join(preprocessed_dir, f"preprocessed_{worker}_{experimentId}.json")
with open(destination_file, "w") as file:
    json.dump(storable, file, indent=4)
    
print(f"Stored file: {destination_file}\n")