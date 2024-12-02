packer {
  required_plugins {
    amazon = {
      version = ">=v1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
    docker = {
      version = ">=v1.0.1"
      source  = "github.com/hashicorp/docker"
    }
  }
}

variable "ami_name_base" {
  type    = string
  default = "terratest-packer-docker-example"
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

data "amazon-ami" "aws" {
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

source "amazon-ebs" "ubuntu-ami" {
  ami_description = "An example of how to create a custom AMI with a simple web app on top of Ubuntu"
  ami_name        = "${var.ami_name_base}-${formatdate("YYYYMMDD-hhmm", timestamp())}"
  encrypt_boot    = false
  instance_type   = var.instance_type
  region          = var.aws_region
  source_ami      = data.amazon-ami.aws.id
  ssh_username    = "ubuntu"
}

source "docker" "ubuntu-docker" {
  changes = ["ENTRYPOINT [\"\"]"]
  commit  = true
  image   = "gruntwork/ubuntu-test:22.04"
}

build {
  sources = ["source.amazon-ebs.ubuntu-ami", "source.docker.ubuntu-docker"]

  provisioner "shell" {
    inline = ["echo 'Sleeping for a few seconds to give Ubuntu time to boot up'", "sleep 30"]
    only   = ["amazon-ebs.ubuntu-ami"]
  }

  provisioner "file" {
    destination = "/tmp/packer-docker-example"
    source      = path.root
  }

  provisioner "shell" {
    inline = ["/tmp/packer-docker-example/configure-sinatra-app.sh"]
  }

  post-processor "docker-tag" {
    only       = ["docker.ubuntu-docker"]
    repository = "gruntwork/packer-docker-example"
    tag        = ["latest"]
  }
}
