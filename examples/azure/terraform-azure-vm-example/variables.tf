# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# ARM_CLIENT_ID
# ARM_CLIENT_SECRET
# ARM_SUBSCRIPTION_ID
# ARM_TENANT_ID

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "disk_type" {
  description = "The type of the Virtual Machine disks"
  type        = string
  default     = "Standard_LRS"
}

variable "location" {
  description = "The Azure location where to deploy your resources too"
  type        = string
  default     = "East US"
}

variable "postfix" {
  description = "A postfix string to centrally mitigate resource name collisions"
  type        = string
  default     = "resource"
}

variable "private_ip" {
  description = "The Static Private IP for the Internal NIC"
  type        = string
  default     = "10.0.17.4"
}

variable "subnet_prefix" {
  description = "The subnet range of IPs for the Virtual Network"
  type        = string
  default     = "10.0.17.0/24"
}

variable "user_name" {
  description = "The username to be provisioned into the vm"
  type        = string
  default     = "testadmin"
}

# Small Windows Server Image, available with Azure Free Account
variable "vm_image_publisher" {
  description = "The storage image reference Publisher from which the VM is created"
  type        = string
  default     = "MicrosoftWindowsServer"
}

variable "vm_image_offer" {
  description = "The storage image reference Offer from which the VM is created"
  type        = string
  default     = "WindowsServer"
}

variable "vm_image_sku" {
  description = "The storage image reference SKU from which the VM is created"
  type        = string
  default     = "2019-Datacenter-Core-smalldisk"
}

variable "vm_image_version" {
  description = "The storage image reference Version from which the VM is created"
  type        = string
  default     = "latest"
}

variable "vm_license_type" {
  description = "The License Type from which the VM is created"
  type        = string
  default     = "Windows_Server"
}

variable "vm_size" {
  description = "The Azure VM Size of the VM"
  type        = string
  default     = "Standard_B1s"
}
