packer {
  required_plugins {
    googlecompute = {
      version = ">=v1.0.0"
      source  = "github.com/hashicorp/googlecompute"
    }
  }
}

variable "gcp_project_id" {
  type    = string
  default = ""
}

variable "gcp_zone" {
  type    = string
  default = "us-central1-a"
}

source "googlecompute" "ubuntu-bionic" {
  image_family        = "terratest"
  image_name          = "terratest-packer-example-${formatdate("YYYYMMDD-hhmm", timestamp())}"
  project_id          = var.gcp_project_id
  source_image_family = "ubuntu-2204-lts"
  ssh_username        = "ubuntu"
  zone                = var.gcp_zone
}


build {
  sources = [
    "source.googlecompute.ubuntu-bionic"
  ]

  provisioner "shell" {
    inline       = ["sudo DEBIAN_FRONTEND=noninteractive apt-get update", "sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y"]
    pause_before = "30s"
  }

}
