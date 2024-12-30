### SUT INSTANCE
resource "google_compute_instance" "fluentd" {
  name = "fluentd-sut"
  machine_type = "e2-highcpu-16"

  boot_disk {
    initialize_params {
	  size = 10
      image = "ubuntu-2004-focal-v20241115"
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