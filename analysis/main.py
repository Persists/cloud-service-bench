import json
import pandas as pd
import matplotlib.pyplot as plt
from datetime import datetime, timedelta
from collections import defaultdict
from tqdm import tqdm
import os

# Directory containing result files
results_dir = "../results"

# Function to estimate the line count in a file
def estimate_line_count(file_path):
    with open(file_path, "r") as file:
        return sum(1 for line in file)

# Initialize a dictionary to store DataFrames for each file
all_dfs = {}

# Process each file in the results directory
for file_name in os.listdir(results_dir):
    if not file_name.endswith(".log"):
        continue
    
    file_path = os.path.join(results_dir, file_name)
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
    cutoff_start = pd.Timedelta(seconds=30)  # Skip the first 30 seconds
    cutoff_end = pd.Timedelta(minutes=8)     # Include only up to 8 minutes after the cutoff

    # Filter the DataFrame to the specified range based on the original timestamps
    df = df[(df.index > start_time + cutoff_start) & (df.index <= start_time + cutoff_start + cutoff_end)]

    # Recalculate the normalized start time to zero seconds
    start_time = df.index[0]
    df['Elapsed Time'] = (df.index - start_time).total_seconds()

    # Convert 'Elapsed Time' to Timedelta and use it as the index for resampling
    df['Elapsed Time'] = pd.to_timedelta(df['Elapsed Time'], unit='s')
    df.set_index('Elapsed Time', inplace=True)

    # Resample the data to 10-second intervals and sum the message counts
    resample_interval = '10s'
    df_resampled = df.resample(resample_interval).sum()
    df_resampled['Throughput per Second'] = df_resampled['Message Count'] / 10

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
