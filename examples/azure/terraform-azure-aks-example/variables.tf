# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "client_id" {
  description = "The Service Principal Client Id for AKS to modify Azure resources."
}
variable "client_secret" {
  description = "The Service Principal Client Password for AKS to modify Azure resources."
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "agent_count" {
  description = "The number of the nodes of the AKS cluster."
  default     = 3
}

variable "ssh_public_key" {
  description = "The public key for the ssh connection to the nodes."
  default     = "~/.ssh/id_rsa.pub"
}

variable "dns_prefix" {
  description = "The prefix to set for the AKS cluster resoureces names."
  default     = "k8stest"
}

variable cluster_name {
  description = "The name to set for the AKS cluster."
  default     = "k8stest"
}

variable resource_group_name {
  description = "The name to set for the resource group."
  default     = "azure-k8stest"
}

variable location {
  description = "The location to set for the AKS cluster."
  default     = "Central US"
}