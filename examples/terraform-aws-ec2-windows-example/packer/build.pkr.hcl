variable "instance_type" {
 type = string
 description = "The EC2 instance size / type to launch"
}

variable "region" {
  type = string
  description = "The AWS region to deploy the Windows instance into"
}


data "amazon-ami" "windows_server_2016" {
  filters = {
    name                = "Windows_Server-2016-English-Full-Base-*"
    root-device-type    = "ebs"
    virtualization-type = "hvm"
  }
  most_recent = true
  owners      = ["801119661308"]
  region      = var.region 
}

locals {
  build_version = "${legacy_isotime("2006.01.02.150405")}"
}

source "amazon-ebs" "windows_server_2016" {
  ami_name                    = "WIN2016-CUSTOM-${local.build_version}"
  associate_public_ip_address = true
  communicator                = "winrm"
  instance_type               = var.instance_type 
  region                      = var.region 
  source_ami                  = "${data.amazon-ami.windows_server_2016.id}"
  user_data_file              = "${path.root}/scripts/bootstrap_windows.txt"
  winrm_timeout               = "15m"
  winrm_password              = "SuperS3cr3t!!!!"
  winrm_username              = "Administrator"

}

build {
  sources = ["source.amazon-ebs.windows_server_2016"]

  # Install Chocolatey package manager, then install any Chocolatey packages defined in scripts/install_packages.ps1
  provisioner "powershell" {
    scripts = ["${path.root}/scripts/install_chocolatey.ps1", "${path.root}/scripts/install_packages.ps1"]
  }

  provisioner "windows-restart" {
    restart_timeout = "35m"
  }
}
