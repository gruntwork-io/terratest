# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_network_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "example" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A VIRTUAL NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "example" {
  name                = "${var.virtual_network_name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.0.0.0/16"]

  tags = {
    environment = "${var.environment_tag}"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY 2 SUBNETS IN THE VIRTUAL NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_subnet" "subnet1" {
  name                      = "${var.first_subnet_name}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  virtual_network_name      = "${azurerm_virtual_network.example.name}"
  address_prefix            = "10.0.11.0/24"
  network_security_group_id = "${azurerm_network_security_group.nsg1.id}"
}

resource "azurerm_subnet" "subnet2" {
  name                      = "${var.second_subnet_name}"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  virtual_network_name      = "${azurerm_virtual_network.example.name}"
  address_prefix            = "10.0.12.0/24"
  network_security_group_id = "${azurerm_network_security_group.nsg2.id}"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A NETWORK SECURITY GROUP FOR EACH SUBNET
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_network_security_group" "nsg1" {
  name                = "${var.first_subnet_nsg_name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

resource "azurerm_subnet_network_security_group_association" "nsg1" {
  subnet_id                 = "${azurerm_subnet.subnet1.id}"
  network_security_group_id = "${azurerm_network_security_group.nsg1.id}"
}

resource "azurerm_network_security_group" "nsg2" {
  name                = "${var.second_subnet_nsg_name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

resource "azurerm_subnet_network_security_group_association" "nsg2" {
  subnet_id                 = "${azurerm_subnet.subnet2.id}"
  network_security_group_id = "${azurerm_network_security_group.nsg2.id}"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A NETWORK SECURITY RULE AND LINK IT TO THE FIRST SUBNET'S NETWORK SECURITY GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_network_security_rule" "example" {
  name                        = "${var.first_subnet_security_rule_name}"
  resource_group_name         = "${azurerm_resource_group.example.name}"
  network_security_group_name = "${azurerm_network_security_group.nsg1.name}"
  priority                    = 111
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A PUBLIC IP ADDRESS RESOURCE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "example" {
  name                = "${var.public_ip_name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  allocation_method   = "Dynamic"
  domain_name_label   = "${var.public_ip_domain_name_label}"
}
