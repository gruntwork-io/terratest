# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE DATA FACTORY
# This is an example of how to deploy an AZURE Data Factory
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------
terraform {
  required_providers {
    azurerm = {
      version = "~> 2.93.0"
      source  = "hashicorp/azurerm"
    }
  }
}
provider "azurerm" {
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE RANDOM PASSWORD
# ---------------------------------------------------------------------------------------------------------------------

# Random password is used as an example to simplify the deployment and improve the security of the database.
# This is not as a production recommendation as the password is stored in the Terraform state file.
resource "random_password" "password" {
  length           = 16
  override_special = "-_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "datafactory_rg" {
  name     = "terratest-datafactory-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A DATA FACTORY
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_data_factory" "data_factory" {
  name                = "datafactory${var.postfix}"
  location            = azurerm_resource_group.datafactory_rg.location
  resource_group_name = azurerm_resource_group.datafactory_rg.name
}