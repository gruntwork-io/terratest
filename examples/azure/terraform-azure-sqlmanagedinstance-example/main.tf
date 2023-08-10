# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE SQL Managed Instance
# This is an example of how to deploy an AZURE SQL Managed Instance
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  version = "~>3.13.0"
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

resource "azurerm_resource_group" "sqlmi_rg" {
  name     = "terratest-sqlmi-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY NETWORK RESOURCES
# This network includes a public address for integration test demonstration purposes
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_network_security_group" "sqlmi_nt_sec_grp" {
  name                = "securitygroup-${var.postfix}"
  location            = azurerm_resource_group.sqlmi_rg.location
  resource_group_name = azurerm_resource_group.sqlmi_rg.name
}

resource "azurerm_network_security_rule" "allow_misubnet_inbound" {
  name                        = "allow_subnet_${var.postfix}"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.0.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.sqlmi_rg.name
  network_security_group_name = azurerm_network_security_group.sqlmi_nt_sec_grp.name
}

resource "azurerm_virtual_network" "sqlmi_vm" {
  name                = "vnet-${var.postfix}"
  resource_group_name = azurerm_resource_group.sqlmi_rg.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.sqlmi_rg.location
}

resource "azurerm_subnet" "sqlmi_sub" {
  name                 = "subnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.sqlmi_rg.name
  virtual_network_name = azurerm_virtual_network.sqlmi_vm.name
  address_prefixes     = ["10.0.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "sqlmi_sb_assoc" {
  subnet_id                 = azurerm_subnet.sqlmi_sub.id
  network_security_group_id = azurerm_network_security_group.sqlmi_nt_sec_grp.id
}

resource "azurerm_route_table" "sqlmi_rt" {
  name                          = "routetable-${var.postfix}"
  location                      = azurerm_resource_group.sqlmi_rg.location
  resource_group_name           = azurerm_resource_group.sqlmi_rg.name
  disable_bgp_route_propagation = false
  depends_on = [
    azurerm_subnet.sqlmi_sub,
  ]
}

resource "azurerm_subnet_route_table_association" "sqlmi_sb_rt_assoc" {
  subnet_id      = azurerm_subnet.sqlmi_sub.id
  route_table_id = azurerm_route_table.sqlmi_rt.id
}

# DEPLOY managed sql instance  ## This depends on vnet ##
resource "azurerm_mssql_managed_instance" "sqlmi_mi" {
  name                = "sqlmi${var.postfix}"
  resource_group_name = azurerm_resource_group.sqlmi_rg.name
  location            = azurerm_resource_group.sqlmi_rg.location

  license_type       = var.sqlmi_license_type
  sku_name           = var.sku_name
  storage_size_in_gb = var.storage_size
  subnet_id          = azurerm_subnet.sqlmi_sub.id
  vcores             = var.cores

  administrator_login          = var.admin_login
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_managed_database" "sqlmi_db" {
  name                = var.sqlmi_db_name
  managed_instance_id = azurerm_mssql_managed_instance.sqlmi_mi.id
}