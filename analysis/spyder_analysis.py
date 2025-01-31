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
# %%

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
        
# %%
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
def read_preprocessed_file(file_path):
    '''This function reads a preprocessed file and returns the data.'''
    with open(file_path, "r") as file:
        return json.load(file)

def read_all_preprocessed_files(directory_path):
    '''This function reads all preprocessed files in a directory and returns the data.'''
    preprocessed_files = []
    for file_name in os.listdir(directory_path):
        if file_name.endswith(".json"):
            preprocessed_files.append(read_preprocessed_file(os.path.join(directory_path, file_name)))
    return preprocessed_files

# %%
# This section was used to preprocess the data from the results directory and store it in a preprocessed directory.
# This was done to speed up the later analysis steps. Reading the data from the log files was the most time-consuming
# part of the analysis process.
directory_path = "../results/1.10worker.memory"
preprocessed_dir = "../results/preprocessed"

worker = 10
experimentId = 10001

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
print(f"Stored file: {destination_file}")

# %%
import pandas as pd
import matplotlib.pyplot as plt

def metrics_to_dataframes(preprocessed_dir):
    '''
    This function reads all preprocessed files in a directory and converts the metrics data to pandas DataFrames.
    '''

    preprocessed_data = read_all_preprocessed_files(preprocessed_dir)
    processed_dateframes = {}

    for data in preprocessed_data:
        worker_count = data["metadata"]["worker"]
        metrics = data["metrics"]
        experimentId = data["metadata"]["expetimentId"]

        df = pd.DataFrame.from_dict(metrics, orient="index")
        df["time"] = pd.to_datetime(df["time"])
        df["used_memory"] = df["mem"].apply(lambda x: x["used"] / (1024 * 1024)) 
        df["cpu"] = df["cpu"].apply(lambda x: sum(x["percentage"]))


        start_time = df["time"].min()
        df["Elapsed Time"] = (df["time"] - start_time).dt.total_seconds() 
        df['Elapsed Time'] = pd.to_timedelta(df['Elapsed Time'], unit='s')

        processed_dateframes[experimentId] = {
            "worker_count": worker_count,
            "df": df
        }

    return processed_dateframes

def plot_cpu(processed_dateframes, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, warm_up_cut=pd.Timedelta(seconds=30), duration_cut=pd.Timedelta(minutes=10)):
    '''
    This function plots the CPU load over time for the given processed dataframes.
    '''

    plotted_workers = set()
    plt.figure(figsize=(12, 6))

    for experimentId, data in processed_dateframes.items():
        worker_count = data["worker_count"]
        if not allow_duplicates and worker in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        df = df[(df["Elapsed Time"] >= warm_up_cut ) & (df["Elapsed Time"] <= duration_cut)]
        label = f"Workers Count: {worker_count}" if not allow_duplicates else f"Worker Count {worker_count} - {experimentId}"
        plt.plot(df["Elapsed Time"] / pd.Timedelta(minutes=1), df["cpu"], label=label)

        plotted_workers.add(worker_count)

    plt.axvspan(0, 30 / 60, color='red', alpha=0.3, label='Warmup Time')
    plt.xlabel("Elapsed Time (minutes)")
    plt.ylabel("CPU Load (%)")
    plt.title("CPU Load Over Time")
    plt.ylim(0, 1600)  # Set y-axis limits
    plt.legend()
    plt.grid(True)
    plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    plt.show()

def plot_memory(processed_dateframes, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, warm_up_cut=pd.Timedelta(seconds=30), duration_cut=pd.Timedelta(minutes=10)):
    '''
    This function plots the memory usage over time for the given processed dataframes.
    '''
    
    plotted_workers = set()
    plt.figure(figsize=(12, 6))

    for experimentId, data in processed_dateframes.items():
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        df = df[(df["Elapsed Time"] >= warm_up_cut ) & (df["Elapsed Time"] <= duration_cut)]
        label = f"Workers Count: {worker_count}" if not allow_duplicates else f"Worker Count {worker_count} - {experimentId}"
        plt.plot(df["Elapsed Time"] / pd.Timedelta(minutes=1), df["used_memory"], label=label)

        plotted_workers.add(worker_count)

    plt.axvspan(0, 30 / 60, color='red', alpha=0.3, label='Warmup Time')
    plt.xlabel("Elapsed Time (minutes)")
    plt.ylabel("Memory Usage (MB)")
    plt.title("Memory Usage Over Time")
    plt.ylim(0, 16*1024)  # Set y-axis limits
    plt.legend()
    plt.grid(True)
    plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    plt.show()

# %%
def throughput_to_dateframes(preprocessed_dir):
    '''
    This function reads all preprocessed files in a directory and converts the throughput data to pandas DataFrames.
    '''

    preprocessed_data = read_all_preprocessed_files(preprocessed_dir)
    processed_dateframes = {}

    for data in preprocessed_data:
        worker_count = data["metadata"]["worker"]
        message_counts = data["message_counts"]
        experimentId = data["metadata"]["expetimentId"]

        df = pd.DataFrame.from_dict(message_counts, orient="index", columns=["message_count"])
        df.index = pd.to_datetime(df.index)
        df.index.name = "time"

        start_time = df.index.min()
        df["Elapsed Time"] = (df.index - start_time).total_seconds()
        df['Elapsed Time'] = pd.to_timedelta(df['Elapsed Time'], unit='s')

        processed_dateframes[experimentId] = {
            "worker_count": worker_count,
            "df": df
        }

    return processed_dateframes

    
# %%

def resample_to_rolling_throughput(throughput_dfs):
    '''This function resamples the throughput data to rolling windows.'''
    for experimentId, data in throughput_dfs.items():
        df = data["df"]
        df = df.sort_index()
        
        # TODO: asks if that is ok
        df["throughput_per_second"] = df["message_count"].rolling(window="4s").sum() / 4        
        throughput_dfs[experimentId]["df"] = df
    
    return throughput_dfs



def plot_throughput_boxplots(throughput_dfs, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True):
    '''
    This function plots boxplots for the throughput data for the given processed dataframes.
    '''
    
    data_to_plot = []
    labels = []
    plotted_workers = set()

    for experimentId, data in throughput_dfs.items():
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        data_to_plot.append(df["throughput_per_second"])
        label = f"Workers Count: {worker_count}" if not allow_duplicates else f"Worker Count {worker_count} - {experimentId}"
        labels.append(label)

        plotted_workers.add(worker_count)

    plt.figure(figsize=(12, 6))
    plt.boxplot(data_to_plot, labels=labels, showfliers=False)
    plt.xlabel("Worker Count")
    plt.ylabel("Throughput (messages/second)")
    plt.title("Throughput Boxplots")
    plt.grid(True)
    plt.show()

def plot_throughput_linechart(throughput_dfs, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, warm_up_cut=pd.Timedelta(seconds=30), duration_cut=pd.Timedelta(minutes=10)):
    '''
    This function plots the throughput over time for the given processed dataframes.
    '''
    
    plotted_workers = set()
    plt.figure(figsize=(12, 6))

    for experimentId, data in throughput_dfs.items():
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        df["throughput_per_second_smothed"] = df["throughput_per_second"].rolling(window=10).mean()
        df = df[(df["Elapsed Time"] >= warm_up_cut ) & (df["Elapsed Time"] <= duration_cut)]
        label = f"Workers Count: {worker_count}" if not allow_duplicates else f"Worker Count {worker_count} - {experimentId}"
        plt.plot(df["Elapsed Time"] / pd.Timedelta(minutes=1), df["throughput_per_second_smothed"], label=label)

        plotted_workers.add(worker_count)

    plt.axvspan(0, 30 / 60, color='red', alpha=0.3, label='Warmup Time')
    plt.xlabel("Elapsed Time (minutes)")
    plt.ylabel("Throughput (messages/second)")
    plt.title("Throughput Over Time")
    plt.legend()
    plt.grid(True)
    plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    plt.show()

# %%
throughput_dfs = throughput_to_dateframes(preprocessed_dir)
resampled_dfs = resample_to_rolling_throughput(throughput_dfs)


plot_throughput_boxplots(resampled_dfs, which_worker_count=[1, 2, 4, 6], allow_duplicates=False)
plot_throughput_linechart(resampled_dfs, which_worker_count=[1, 2, 4, 6], allow_duplicates=False)
#plot_throughput_boxplots(resampled_dfs, which_worker_count=[8, 10, 12])

# %%
preprocessed_dir="./results/preprocessed/"
dfs=metrics_to_dataframes(preprocessed_dir)

print("test")
plot_cpu(dfs, which_worker_count=[1, 2, 4, 6], allow_duplicates=False)
plot_cpu(dfs, which_worker_count=[8,10,12])
plot_memory(dfs, which_worker_count=[6], allow_duplicates=True)


