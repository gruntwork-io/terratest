# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# ARM_TENANT_ID
# ARM_SUBSCRIPTION_ID
# ARM_CLIENT_ID (if you are using a Service Principal to authenticate to Azure)
# ARM_CLIENT_SECRET (if you are using a Service Principal to authenticate to Azure)

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "location" {
  description = "The Azure region where resources will be located."
}

variable "public_ip_domain_name_label" {
  description = "The domain name label (DNS prefix) to set for the public IP."
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "resource_group_name" {
  description = "The name of the resource group where the resources will be deployed."
  default     = "terratest-example-rg"
}

variable "virtual_network_name" {
  description = "The name to set for the virtual network."
  default     = "terratest-vnet"
}

variable "environment_tag" {
  description = "Value to set for the 'environment' tag applied to the virtual network."
  default     = "test"
}

variable "first_subnet_name" {
  description = "The name to set for the first subnet."
  default     = "terratest-subnet1"
}

variable "second_subnet_name" {
  description = "The name to set for the second subnet."
  default     = "terratest-subnet2"
}

variable "first_subnet_nsg_name" {
  description = "The name of network security group to apply to the first subnet."
  default     = "subnet1-nsg"
}

variable "second_subnet_nsg_name" {
  description = "The name of network security group to apply to the second subnet."
  default     = "subnet2-nsg"
}

variable "first_subnet_security_rule_name" {
  description = "The name of the network security rule to link to the first subnet's network security group."
  default     = "Allow_SSH_Inbound"
}

variable "public_ip_name" {
  description = "The name to set for the public IP resource."
  default     = "terratest-example-ip"
}
