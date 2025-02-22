resource "google_compute_instance_template" "go_client_template" {  
    machine_type = "e2-standard-4"
    disk {
        source_image = "ubuntu-2004-focal-v20241115"
        disk_size_gb = 25
    }
    metadata_startup_script = file("startup_go_client.sh")
    network_interface {
        network = google_compute_network.vpc_network.id
        access_config {
          # Include this section to give the VM an external IP address
        }
    }
}

resource "google_compute_instance_from_template" "generator" {
  count                   = var.generator_count
  name                    = format("generator-%02d", count.index + 1)
  source_instance_template = google_compute_instance_template.go_client_template.id
}

resource "google_compute_instance" "sink" {
  count = var.sink_count
  name = format("sink-%02d", count.index + 1)
  machine_type = "e2-standard-16"

  boot_disk {
      initialize_params {
        image        = "ubuntu-2004-focal-v20241115"
        size         = 50
      }
  }

  metadata_startup_script = file("startup_go_client.sh")

  network_interface {
    network = google_compute_network.vpc_network.id
    access_config {
      # Include this section to give the VM an external IP address
    }
  }
}
