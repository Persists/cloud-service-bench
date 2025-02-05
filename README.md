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
- `analysis/`: Contains a Spyder notebook to preprocess and analyze the data.
- `scripts/`: Contains scripts to run the benchmark and load the resulting data from the GCP instance.

## Setup

### Creating the VMs:

To deploy the infrastructure on GCP, follow these steps:

In the Root directory of this repository, run the following commands:

- The following command will create 4 load generators and 1 sink,
  we used that configuration to run the benchmark

```bash
terraform -chdir=infrastructure apply -var="generator_count=4" -var="sink_count=1"
```

This command will create the necessary infrastructure on GCP.

### Checking the Benchmark VMs:

We use metadata_startup_script's to install the necessary software on the VMs. To check if the software is installed, the following script can be run:

- The following command will check if the software is installed on the VMs, the first argument is the number of generators and the second argument is the number of sinks.

```bash
sudo bash ./scripts/check_startup_script.sh 4 1
```

### Start the SUT (System Under Test) VM:

To start the fluentd instance, run the following command:

- The following command will start the fluentd instance on the sink VM.
- It will use the fluentd configuration file `.config/fluentd/fluentd.conf`

```bash
sudo bash ./scripts/start_sut.sh
```

### Running the Benchmark:

To run the benchmark, follow these steps:

- The following command will start the load generators, the sinks and the monitor on the respective VMs.
- The first argument is the number of generators and the second argument is the number of sinks.
- The configuration of the benchmark can be found in the `.config/experiment/config.yaml` file.

```bash
sudo bash ./scripts/start_benchmark.sh 4 1
```

### End of the Benchmark:

The benchmark will be stopped automatically after the duration specified in the configuration file.

### Getting the Data:

To get the data from the VMs, follow these steps:

- The following command will get the data from the VMs and save it in the `data` directory.
- Finally it will destroy the VMs.
- The first argument is the number of generators and the second argument is the number of sinks.
- The last argument is the experiment id.

```bash
sudo bash ./scripts/get_data.sh 4 1 1001
```

### Analyzing the Data:

To analyze the data, use the spyder notebook in the `analysis` directory.

- First use the scriots in `preprocess.py` to preprocess the data.
- Then use the `analysis.py` to analyze the data.

### Helpfull commands:

During the benchmark, htop can be used to monitor the resource utilization of the VMs.

Destroying only the SUT instance:

```bash
terraform -chdir=infrastructure destroy -auto-approve -target="google_compute_instance.fluentd"
```

Destroying the generator instances:

- The following command will destroy the generator instances, it destroys the instances from 0 to 3.

```bash
terraform -chdir=infrastructure destroy -auto-approve $(printf -- '-target=google_compute_instance_from_template.generator[%d] ' {0..3})
```

Destroying everything:

```bash
terraform -chdir=infrastructure destroy -auto-approve
```
