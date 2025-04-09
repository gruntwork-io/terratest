packer {
  required_plugins {
    amazon = {
      version = ">=v1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
    oracle = {
      version = ">=v1.0.0"
      source  = "github.com/hashicorp/oracle"
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

data "amazon-ami" "ubuntu-jammy" {
  filters = {
    architecture                       = "x86_64"
    "block-device-mapping.volume-type" = "gp2"
    name                               = "*ubuntu-jammy-22.04-amd64-server-*"
    root-device-type                   = "ebs"
    virtualization-type                = "hvm"
  }
  most_recent = true
  owners      = ["099720109477"]
  region      = var.aws_region
}

source "amazon-ebs" "ubuntu-example" {
  ami_description = "An example of how to create a custom AMI on top of Ubuntu"
  ami_name        = "${var.ami_base_name}-terratest-packer-example"
  encrypt_boot    = false
  instance_type   = var.instance_type
  region          = var.aws_region
  source_ami      = data.amazon-ami.ubuntu-jammy.id
  ssh_username    = "ubuntu"
}

source "oracle-oci" "oracle-example" {
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
    "source.amazon-ebs.ubuntu-example",
    "source.oracle-oci.oracle-example"
  ]

  provisioner "shell" {
    inline       = ["sudo DEBIAN_FRONTEND=noninteractive apt-get update", "sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y"]
    pause_before = "30s"
  }
}
