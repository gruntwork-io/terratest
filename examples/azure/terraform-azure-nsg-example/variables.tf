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

variable "postfix" {
  description = "Random postfix string used for each test run; set from the test file at runtime."
  type        = string
  default     = "qwefgt"
}

variable "resource_group_name" {
  description = "Name for the resource group holding resources for this example"
  type        = string
  default     = "terratest-nsg-rg"
}

variable "location" {
  description = "The Azure region in which to deploy this sample"
  type        = string
  default     = "East US"
}

variable "vnet_name" {
  description = "Name for the example virtual network"
  type        = string
  default     = "vnet01"
}

variable "subnet_name" {
  description = "Name for the example virtual network default subnet"
  type        = string
  default     = "subnet01"
}

variable "vm_nic_name" {
  description = "Name for the NIC attached to our example VM."
  type        = string
  default     = "nic01"
}

variable "vm_nic_ip_config_name" {
  description = "Name for the NIC IP configuration attached to our example VM."
  type        = string
  default     = "nic_ipconfig01"
}

variable "nsg_name" {
  description = "Name for the example NSG."
  type        = string
  default     = "nsg01"
}

variable "nsg_ssh_rule_name" {
  description = "Name for the example SSH NSG rule used in this example."
  type        = string
  default     = "nsgrule01"
}

variable "nsg_http_rule_name" {
  description = "Name for the example HTTP NSG rule used in this example."
  type        = string
  default     = "nsgrule02"
}

variable "vm_name" {
  description = "The name of the VM used in this example"
  type        = string
  default     = "vm01"
}

variable "vm_size" {
  description = "The size of the VM to deploy"
  type        = string
  default     = "Standard_D2s_v3"
}

variable "hostname" {
  description = "The hostname of the new VM to be configured"
  type        = string
  default     = "vm01"
}

variable "os_disk_name" {
  description = "The of the OS disk to use on our example VM."
  type        = string
  default     = "osdisk01"
}

variable "username" {
  description = "The username to be provisioned into your VM"
  type        = string
  default     = "testadmin"
}

variable "password" {
  description = "The password to configure for SSH access"
  type        = string
  default     = "!@#PasswordSetInCode!@#"
}
