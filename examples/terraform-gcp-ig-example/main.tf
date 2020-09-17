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
# DEPLOY A REGIONAL MANAGED INSTANCE GROUP
# See test/terraform_gcp_ig_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

# Create a Regional Managed Instance Group
resource "google_compute_region_instance_group_manager" "example" {
  project = var.gcp_project_id
  region  = var.gcp_region

  name               = "${var.cluster_name}-ig"
  base_instance_name = var.cluster_name
  version {
    name              = "terratest"
    instance_template = google_compute_instance_template.example.self_link
  }

  target_size = var.cluster_size
}

# Create the Instance Template that will be used to populate the Managed Instance Group.
resource "google_compute_instance_template" "example" {
  project = var.gcp_project_id

  name_prefix  = var.cluster_name
  machine_type = var.machine_type

  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
    preemptible         = false
  }

  disk {
    boot         = true
    auto_delete  = true
    source_image = "ubuntu-os-cloud/ubuntu-1604-lts"
  }

  network_interface {
    network = "default"

    # The presence of this property assigns a public IP address to each Compute Instance. We intentionally leave it
    # blank so that an external IP address is selected automatically.
    access_config {
    }
  }
}

