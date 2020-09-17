# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A CLOUD INSTANCE RUNNING UBUNTU
# See test/terraform_gcp_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "google_compute_instance" "example" {
  project      = var.gcp_project_id
  name         = var.instance_name
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1604-lts"
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A GOOGLE STORAGE BUCKET
# ---------------------------------------------------------------------------------------------------------------------

resource "google_storage_bucket" "example_bucket" {
  project  = var.gcp_project_id
  name     = var.bucket_name
  location = var.bucket_location
}

