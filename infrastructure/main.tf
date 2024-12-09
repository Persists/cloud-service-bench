provider "google" {
    project = "cloud-service-be"
    region = "europe-west3"
    zone = "europe-west3-c"
}

resource "google_compute_network" "vpc_network" {
    name = "vpc-network"
    auto_create_subnetworks = true
}

### FIREWALL
resource "google_compute_firewall" "all" {
  name = "allow-all"
  allow {
    protocol = "tcp"
    ports = ["0-65535"]
  }
  network = google_compute_network.vpc_network.id
  source_ranges = ["0.0.0.0/0"]
}

### SUT INSTANCE
resource "google_compute_instance" "victoria" {
  name = "victoria-r1"
  machine_type = "e2-standard-2"

  boot_disk {
    initialize_params {
	  size = 10
      image = "ubuntu-2004-focal-v20231101"
    }
  }
    metadata_startup_script = file("startup_sut.sh")


  network_interface {
    network = google_compute_network.vpc_network.id
    access_config {
      # Include this section to give the VM an external IP address
    }
  }
}

