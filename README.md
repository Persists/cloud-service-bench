# cloud-service-bench

This repository contains the source code for a cloud service benchmarking tool. The tool is designed to benchmark the performance of fluentd, a popular log collector, on Google Cloud Platform (GCP).

## Structure

The repository is structured as follows:

- `cmd/`: Contains the main.go files for all necessary components, which are used to run the benchmark.

  - `monitor/`: Monitor resource utilization of the fluentd instance.
  - `http-sink/`: Receive the logs from fluentd.
  - `load-generator/`: Generate logs and send them to fluentd.

- `internal/`: Contains the internal packages used by the components.
- `infrastucture/`: Contains the terraform files to create the necessary infrastructure on GCP.
- `analysis/`: Contains a Spyder notebook to analyze the results of the benchmark.
- `scripts/`: Contains scripts to run the benchmark and load the resulting data from the GCP instance.

## Setup

### Creating the VMs:

To deploy the infrastructure on GCP, follow these steps:
