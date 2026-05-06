terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  region = "us-east1"
}

# website::tag::1:: Deploy a cloud instance
resource "google_compute_instance" "example" {
  name         = var.instance_name
  machine_type = "f1-micro"
  zone         = "us-east1-b"

  # website::tag::2:: Run Ubuntu 22.04 on the instance
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }
}

# website::tag::3:: Allow the user to pass in a custom name for the instance
variable "instance_name" {
  description = "The Name to use for the Cloud Instance."
  default     = "gcp-hello-world-example"
}
