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