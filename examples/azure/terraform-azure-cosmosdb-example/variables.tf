# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "resource_group_name" {
  description = "The name to set for the resource group."
  default     = "azure-cosmosdb-test"
}

variable "location" {
  description = "The location to set for the CosmosDB instance."
  default     = "East US"
}

variable "cosmosdb_account_name" {
  description = "The name to set for the CosmosDB account."
  default     = "azure-cosmosdb"
}
