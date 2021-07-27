packer {
  required_version = ">= 1.7.0"
  required_plugins {
    # This plugin is primarily used for testing the new required_plugins feature introduced in Packer 1.7.0.
    comment = {
        version = ">=v0.2.23"
        source = "github.com/sylviamoss/comment"
    }
  }
}

variable "ami_base_name" {
  type    = string
  default = ""
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "gcp_project_id" {
  type    = string
  default = ""
}

variable "gcp_zone" {
  type    = string
  default = "us-central1-a"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

variable "oci_availability_domain" {
  type    = string
  default = ""
}

variable "oci_base_image_ocid" {
  type    = string
  default = ""
}

variable "oci_compartment_ocid" {
  type    = string
  default = ""
}

variable "oci_pass_phrase" {
  type    = string
  default = ""
}

variable "oci_subnet_ocid" {
  type    = string
  default = ""
}

data "amazon-ami" "ubuntu-xenial" {
  filters = {
    architecture                       = "x86_64"
    "block-device-mapping.volume-type" = "gp2"
    name                               = "*ubuntu-xenial-16.04-amd64-server-*"
    root-device-type                   = "ebs"
    virtualization-type                = "hvm"
  }
  most_recent = true
  owners      = ["099720109477"]
  region      = var.aws_region
}

source "amazon-ebs" "aws" {
  ami_description = "An example of how to create a custom AMI on top of Ubuntu"
  ami_name        = "${var.ami_base_name}-terratest-packer-example"
  encrypt_boot    = false
  instance_type   = var.instance_type
  region          = var.aws_region
  source_ami      = data.amazon-ami.ubuntu-xenial.id
  ssh_username    = "ubuntu"
}

source "googlecompute" "gcp" {
  image_family        = "terratest"
  image_name            = "terratest-packer-example-${formatdate("YYYYMMDD-hhmm", timestamp())}"
  project_id          = var.gcp_project_id
  source_image_family = "ubuntu-1804-lts"
  ssh_username        = "ubuntu"
  zone                = var.gcp_zone
}

source "oracle-oci" "oracle" {
  availability_domain = var.oci_availability_domain
  base_image_ocid     = var.oci_base_image_ocid
  compartment_ocid    = var.oci_compartment_ocid
  image_name          = "terratest-packer-example-${formatdate("YYYYMMDD-hhmm", timestamp())}"
  pass_phrase         = var.oci_pass_phrase
  shape               = "VM.Standard2.1"
  ssh_username        = "ubuntu"
  subnet_ocid         = var.oci_subnet_ocid
}

build {
  sources = [
    "source.amazon-ebs.aws",
    "source.googlecompute.gcp",
    "source.oracle-oci.oracle"
  ]

  provisioner "comment" {
    comment     = "Basic comment example"
    ui          = false
  }

  provisioner "shell" {
    inline       = ["sudo DEBIAN_FRONTEND=noninteractive apt-get update", "sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y"]
    pause_before = "30s"
  }

}