import json
import pandas as pd
import matplotlib.pyplot as plt
from datetime import datetime
from tqdm import tqdm
from collections import defaultdict
import os

# Path configurations
CPU_LOG_FILE = "../results/1.2worker.memory/data/fluentd-sut/monitor_fluentd-sut_3_w16.log"
RESULTS_DIR = ["../results/1.4worker.memory/data/sink-01", "../results/2.4worker.memory/data/sink-01", "../results/3.4worker.memory/data/sink-01"]

# Plot CPU Load
def plot_cpu_load(cpu_log_file):
    timestamps = []
    cpu_loads = []

    try:
        with open(cpu_log_file, "r") as file:
            for line in file:
                # Skip metadata or non-JSON lines
                if not line.strip().startswith("{"):
                    continue
                # Parse the JSON line
                data = json.loads(line)
                timestamp_str = data["time"]
                timestamp = datetime.fromisoformat(timestamp_str.replace("Z", "+00:00"))
                cpu_load = data["cpu"]["percentage"]
                timestamps.append(timestamp)
                cpu_loads.append(cpu_load)

        # Convert to DataFrame
        df = pd.DataFrame(cpu_loads, index=timestamps)
        df.index.name = "Timestamp"
        df.columns = [f"CPU {i}" for i in range(len(df.columns))]

        # Calculate total and average CPU load
        df["Total CPU"] = df.sum(axis=1)
        df["Average CPU"] = df["Total CPU"] / len(df.columns)

        # Plot the CPU load
        plt.figure(figsize=(12, 6))
        for column in df.columns[:-2]:  # Exclude "Total CPU" and "Average CPU" columns
            plt.plot(df.index, df[column], alpha=0.3, label=column)
        plt.plot(df.index, df["Total CPU"], color="blue", label="Total CPU Load (1600%)", linewidth=2)
        plt.plot(df.index, df["Average CPU"], color="red", label="Average CPU Load (%)", linewidth=2)

        # Configure the plot
        plt.title("CPU Load Over Time (Scaled to 1600%)")
        plt.xlabel("Time")
        plt.ylabel("CPU Load (%)")
        plt.legend(loc="upper left", fontsize="small")
        plt.grid(True)
        plt.tight_layout()

        # Show the plot
        plt.show()

    except FileNotFoundError:
        print(f"File not found: {cpu_log_file}")
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
    except Exception as e:
        print(f"An error occurred: {e}")

# Plot Throughput
def plot_throughput(results_dirs):
    all_dfs = {}

    def estimate_line_count(file_path):
        with open(file_path, "r") as file:
            return sum(1 for line in file)

    # Process each directory or file in the results_dirs list
    for dir_or_file in results_dirs:
        if os.path.isfile(dir_or_file):  # If it's a single file
            files_to_process = [dir_or_file]
        else:  # If it's a directory
            files_to_process = [os.path.join(dir_or_file, f) for f in os.listdir(dir_or_file) if f.endswith(".log")]

        for file_path in files_to_process:
            file_name = os.path.basename(file_path)
            message_counts = defaultdict(int)
            experiment_started = False
            experiment_ended = False

            progress_bar = tqdm(total=estimate_line_count(file_path), desc=f"Processing {file_name}", unit="line", dynamic_ncols=True)

            try:
                with open(file_path, "r") as file:
                    for line in file:
                        progress_bar.update(1)
                        if "Experiment phase" in line:
                            experiment_started = True
                            continue

                        if "Cool-down phase" in line:
                            experiment_ended = True
                            progress_bar.n = progress_bar.total
                            progress_bar.refresh()
                            break

                        if experiment_started and not experiment_ended:
                            try:
                                data = json.loads(line)
                            except json.JSONDecodeError:
                                continue

                            timestamp_str = data["arrivalTime"]
                            timestamp = datetime.fromisoformat(timestamp_str)
                            message_count = len(data.get("logMessages", []))
                            message_counts[timestamp] += message_count
            except json.JSONDecodeError as e:
                print(f"Error decoding JSON: {e}")

            # Convert the message counts dictionary into a DataFrame
            df = pd.DataFrame(list(message_counts.items()), columns=['Timestamp', 'Message Count'])
            df['Timestamp'] = pd.to_datetime(df['Timestamp'])
            df = df.sort_values(by='Timestamp')
            df.set_index('Timestamp', inplace=True)

            # Normalize the start time to zero seconds
            start_time = df.index[0]

            # Define time limits for analysis
            cutoff_start = pd.Timedelta(seconds=0)  # Skip the first 0 seconds
            cutoff_end = pd.Timedelta(minutes=8)    # Include only up to 8 minutes after the cutoff

            # Filter the DataFrame to the specified range based on the original timestamps
            df = df[(df.index > start_time + cutoff_start) & (df.index <= start_time + cutoff_start + cutoff_end)]

            # Recalculate the normalized start time to zero seconds
            start_time = df.index[0]
            df['Elapsed Time'] = (df.index - start_time).total_seconds()

            # Convert 'Elapsed Time' to Timedelta and use it as the index for resampling
            df['Elapsed Time'] = pd.to_timedelta(df['Elapsed Time'], unit='s')
            df.set_index('Elapsed Time', inplace=True)

            # Resample the data to 1-second intervals and sum the message counts
            resample_interval = '1s'
            df_resampled = df.resample(resample_interval).sum()
            df_resampled['Throughput per Second'] = df_resampled['Message Count'] / 1

            all_dfs[file_name] = df_resampled
            progress_bar.close()

    # Plot all data on one plot
    plt.figure(figsize=(15, 7))

    for file_name, df_resampled in all_dfs.items():
        plt.plot(df_resampled.index.total_seconds() / 60, df_resampled['Throughput per Second'], label=file_name)

    plt.xticks(rotation=45, ha='right')
    plt.xlabel('Elapsed Time (minutes)')
    plt.ylabel('Throughput per Second (Messages per Second)')
    plt.title('Normalized Throughput per Second for All Files')
    plt.legend()
    plt.tight_layout()
    plt.show()


# Run both plots
# plot_cpu_load(CPU_LOG_FILE)
# plot_throughput(RESULTS_DIR)
# plot_throughput(["../results/1.6worker.memory/data/sink-01", "../results/2.6worker.memory/data/sink-01", "../results/4.6worker.memory/data/sink-01"])
#plot_throughput(["../results/1.2worker.memory/data/sink-01", "../results/2.2worker.memory/data/sink-01", "../results/3.2worker.memory/data/sink-01"])
plot_throughput(["../results/1.1worker.memory/data/sink-01", "../results/2.1worker.memory/data/sink-01", "../results/3.1worker.memory/data/sink-01"])