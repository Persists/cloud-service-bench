#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Sat Jan 25 21:00:27 2025

@author: jonasheisterberg
"""

# %%
import json
import os
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
        
preprocessed_dir = "./results/preprocessed"
warm_up_cut=pd.Timedelta(seconds=30)
duration_cut=pd.Timedelta(minutes=10)

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
        df = df[(df["Elapsed Time"] >= warm_up_cut ) & (df["Elapsed Time"] <= duration_cut)]


        processed_dateframes[experimentId] = {
            "worker_count": worker_count,
            "df": df
        }

    return processed_dateframes

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

    
def resample_to_rolling_throughput(throughput_dfs):
    '''This function resamples the throughput data to rolling windows.'''
    for experimentId, data in throughput_dfs.items():
        df = data["df"]
        df = df.sort_index()
        
        df["throughput_per_second"] = df["message_count"].rolling(window="4s").sum() / 4
        # this is the smoothed throughput, it is only used for the line chart, for better visualization
        df["throughput_per_second_smoothed"] = df["throughput_per_second"].rolling(window=10).mean()

        df = df[(df["Elapsed Time"] >= warm_up_cut ) & (df["Elapsed Time"] <= duration_cut)]

        throughput_dfs[experimentId]["df"] = df
    
    return throughput_dfs

def plot_cpu(processed_dateframes, title, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, which_id=None):
    '''
    This function plots the CPU load over time for the given processed dataframes.
    '''

    plotted_workers = set()
    plt.figure(figsize=(12, 6))

    for experimentId, data in processed_dateframes.items():
        if which_id is not None and experimentId not in which_id:
            continue
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        label = f"Number of Workers: {worker_count}" if not allow_duplicates else f"Number of Workers: {worker_count} \n (id: {experimentId})"
        plt.plot(df["Elapsed Time"] / pd.Timedelta(minutes=1), df["cpu"], label=label)

        plotted_workers.add(worker_count)

    plt.axvspan(0, 30 / 60, color='red', alpha=0.3, label='Warmup Time')
    plt.xlabel("Elapsed Time (minutes)")
    plt.ylabel("CPU Load (%)")
    plt.title(title)
    plt.ylim(0, 1600)  # Set y-axis limits
    plt.legend()
    plt.grid(True)
    plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    plt.show()

def plot_memory(processed_dateframes, title, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, which_id=None):
    '''
    This function plots the memory usage over time for the given processed dataframes.
    '''
    
    plotted_workers = set()
    plt.figure(figsize=(12, 6))

    for experimentId, data in processed_dateframes.items():
        if which_id is not None and experimentId not in which_id:
            continue
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        label = f"Number of Workers: {worker_count}" if not allow_duplicates else f"Number of Workers: {worker_count} \n (id: {experimentId})"
        plt.plot(df["Elapsed Time"] / pd.Timedelta(minutes=1), df["used_memory"], label=label)

        plotted_workers.add(worker_count)

    plt.axvspan(0, 30 / 60, color='red', alpha=0.3, label='Warmup Time')
    plt.xlabel("Elapsed Time (minutes)")
    plt.ylabel("Memory Usage (MB)")
    plt.title(title)
    plt.ylim(0, 4*1024)  # Set y-axis limits
    plt.legend()
    plt.grid(True)
    plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    plt.show()

def plot_throughput_boxplots(throughput_dfs, title, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, which_id=None):
    '''
    This function plots boxplots for the throughput data for the given processed dataframes and prints mean and quantiles.
    '''
    
    data_to_plot = []
    labels = []
    plotted_workers = set()

    for experimentId, data in throughput_dfs.items():
        if which_id is not None and experimentId not in which_id:
            continue
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        throughput_values = df["throughput_per_second"].dropna()
        data_to_plot.append(throughput_values)

        label = f"{worker_count}" if not allow_duplicates else  f"Number of Workers: {worker_count} \n (id: {experimentId})"
        labels.append(label)

        plotted_workers.add(worker_count)

    plt.figure(figsize=(12, 6))
    plt.boxplot(data_to_plot, labels=labels, showfliers=False)
    plt.ylabel("Throughput (logs/second)")
    plt.title(title)
    plt.grid(True)
    plt.show()





def plot_throughput_linechart(throughput_dfs, title, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, which_id=None):
    '''
    This function plots the throughput over time for the given processed dataframes.

    It uses the smoothed throughput for better visualization.
    '''
    
    plotted_workers = set()
    plt.figure(figsize=(12, 6))

    for experimentId, data in throughput_dfs.items():
        if which_id is not None and experimentId not in which_id:
            continue
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        df = data["df"]
        label = f"Number of Workers: {worker_count}" if not allow_duplicates else f"Number of Workers: {worker_count} \n (id: {experimentId})"

        plt.plot(df["Elapsed Time"] / pd.Timedelta(minutes=1), df["throughput_per_second_smoothed"], label=label)

        plotted_workers.add(worker_count)

    plt.axvspan(0, 30 / 60, color='red', alpha=0.3, label='Warmup Time')
    plt.xlabel("Elapsed Time (minutes)")
    plt.ylabel("Throughput (logs/second)")
    plt.title(title)
    plt.legend()    
    plt.grid(True)
    plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    plt.show()

    # %%

def plot_throughput_barplots(cleaned_dfs, title, which_worker_count=[1, 2, 4, 6, 8, 10, 12], allow_duplicates=True, which_id=None):
    '''
    This function plots barplots for the throughput data for the given processed dataframes.

    It uses the mean throughput for each experiment.
    '''
    
    data_to_plot = []
    labels = []
    plotted_workers = set()
    colors = []
    ids = []
    # reverse bars printed

    cleaned_dfs = dict(sorted(cleaned_dfs.items(), key=lambda x: x[0], reverse=False))


    for experimentId, data in cleaned_dfs.items():
        if which_id is not None and experimentId not in which_id:
            continue
        worker_count = data["worker_count"]
        if not allow_duplicates and worker_count in plotted_workers or worker_count not in which_worker_count:
            continue

        data_to_plot.append(data["stats"]["mean"])


        label = f"{experimentId}" if not allow_duplicates else f"{experimentId}"
        labels.append(label)
        colors.append(worker_count)
        ids.append(experimentId)

        plotted_workers.add(worker_count)

    unique_worker_counts = sorted(set(colors))
    color_map = {worker_count: plt.cm.tab20(i / len(unique_worker_counts)) for i, worker_count in enumerate(unique_worker_counts)}
    bar_colors = [color_map[worker_count] for worker_count in colors]

    plt.figure(figsize=(14, 6))
    bars = plt.bar(labels, data_to_plot, color=bar_colors)
    plt.ylabel("Throughput (logs/second)")
    plt.xlabel("Experiment ID")
    plt.title(title)
    plt.grid(True)

    for bar, mean in zip(bars, data_to_plot):
        yval = bar.get_height()
        plt.text(bar.get_x() + bar.get_width()/2, yval, round(mean), va='bottom', ha='center')

    handles = [plt.Line2D([0], [0], color=color_map[worker_count], lw=4) for worker_count in unique_worker_counts]
    plt.legend(handles, [f"Workers: {wc}" for wc in unique_worker_counts], title="Number of Workers", bbox_to_anchor=(1.05, 1), loc='upper left')

    plt.show()


def clean_the_throughput_data(throughput_dfs):
    '''
    This function cleans the throughput data by removing outliers.

    The function uses the IQR method to remove outliers.

    The function returns a dictionary with the cleaned dataframes and the statistics for each dataframe.
    '''
    cleaned_throughput_dfs = {}
    for experimentId, data in throughput_dfs.items():
        df = data["df"]

        Q1 = df["throughput_per_second"].quantile(0.25)
        Q3 = df["throughput_per_second"].quantile(0.75)
        IQR = Q3 - Q1

        lower_bound = Q1 - 1.5 * IQR
        upper_bound = Q3 + 1.5 * IQR

        cleaned_df = df[(df["throughput_per_second"] >= lower_bound) & (df["throughput_per_second"] <= upper_bound)]


        min = cleaned_df["throughput_per_second"].min()
        max = cleaned_df["throughput_per_second"].max()

        median = cleaned_df["throughput_per_second"].median()
        mean = cleaned_df["throughput_per_second"].mean()
        std = cleaned_df["throughput_per_second"].std()

        
        cleaned_throughput_dfs[experimentId] = {
            "worker_count": data["worker_count"],
            "df": cleaned_df,
            "stats": {
                "min": min,
                "max": max,
                "mean": mean,
                "median": median,
                "std": std,
                "Q1": Q1,
                "Q3": Q3,

            }
        }

    return cleaned_throughput_dfs

# %%
# Read the preprocessed data and convert it to pandas DataFrames

throughput_dfs = throughput_to_dateframes(preprocessed_dir)
resampled_raw_dfs = resample_to_rolling_throughput(throughput_dfs)

metrics_dfs=metrics_to_dataframes(preprocessed_dir)

# sort and reverse dfs by experiment id
resampled_raw_dfs = dict(sorted(resampled_raw_dfs.items(), key=lambda x: x[0], reverse=True))
metrics_dfs = dict(sorted(metrics_dfs.items(), key=lambda x: x[0], reverse=True))

# %%
# Plot the raw data

plot_throughput_boxplots(resampled_raw_dfs, "Throughput Fluentd with Different Number of Workers",
                         allow_duplicates=False, which_id=[
                            1001, 2001, 4001, 6001
                        ])

plot_throughput_linechart(resampled_raw_dfs, "Throughput Fluentd with Different Number of Workers",
                         allow_duplicates=False, which_id=[
                            1001, 2001, 4001, 6001
                        ])

plot_cpu(metrics_dfs, "CPU Load Fluentd with Different Number of Workers",
         allow_duplicates=False, which_id=[
            1001, 2001, 4001, 6001
        ])

plot_memory(metrics_dfs, "Memory Usage Fluentd with Different Number of Workers",
            allow_duplicates=False, which_id=[
                1001, 2002, 4002, 6002
            ])

plot_throughput_boxplots(resampled_raw_dfs, "Throughput Fluentd for 6 Concurrent Workers",
                         allow_duplicates=True, which_worker_count=[6])

plot_throughput_boxplots(resampled_raw_dfs, "Throughput Fluentd for 4 Concurrent Workers",
                            allow_duplicates=True, which_worker_count=[4])

plot_throughput_boxplots(resampled_raw_dfs, "Throughput Fluentd for 2 Concurrent Workers",
                            allow_duplicates=True, which_worker_count=[2])

plot_throughput_boxplots(resampled_raw_dfs, "Throughput Fluentd for 1 Concurrent Workers",
                            allow_duplicates=True, which_worker_count=[1])


# %%
# Plot the cleaned data (without outliers)
cleaned_throughput_dfs = clean_the_throughput_data(resampled_raw_dfs)

plot_throughput_barplots(cleaned_throughput_dfs, "Mean Throughput for Different Number of Workers",
                            which_worker_count=[1,2,4,6])


# save stats to csv
stats = []
for experimentId, data in cleaned_throughput_dfs.items():
    stats.append({
        "experimentId": experimentId,
        "worker_count": data["worker_count"],
        "min": data["stats"]["min"],
        "max": data["stats"]["max"],
        "median": data["stats"]["median"],
        "mean": data["stats"]["mean"],
        "std": data["stats"]["std"],
        "Q1": data["stats"]["Q1"],
        "Q3": data["stats"]["Q3"],
    })


stats_df = pd.DataFrame(stats)
stats_df.to_csv("results/stats.csv", index=False)


# %%
# Plot the anomalies

plot_throughput_linechart(resampled_raw_dfs, "Anomalies in Throughput Fluentd with Different Number of Workers",
                             which_worker_count=[8, 10, 12])

plot_throughput_barplots(cleaned_throughput_dfs, "Anomalies in Mean Throughput for Different Number of Workers",
                            which_worker_count=[1,2,4,6, 8, 10, 12])

plot_cpu(metrics_dfs, "Anomalies in CPU Load Fluentd with 8 Workers",
            which_worker_count=[8])

plot_cpu(metrics_dfs, "Anomalies in CPU Load Fluentd with 10 and 12 Workers",
            which_worker_count=[10, 12])

plot_memory(metrics_dfs, "Anomalies in Memory Usage Fluentd with 8 Workers",
                which_worker_count=[8])

plot_memory(metrics_dfs, "Anomalies in Memory Usage Fluentd with Different Number of Workers",
                which_worker_count=[8, 10, 12])

plot_throughput_boxplots(resampled_raw_dfs, "Anomalies in Throughput Fluentd with 8 Workers",
                            which_worker_count=[8])

plot_throughput_boxplots(resampled_raw_dfs, "Anomalies in Throughput Fluentd with 10 and 12 Workers",
                            which_worker_count=[10, 12])