# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE AUTOMATION ACCOUNT, AN AZURE VM ALONG WITH AN EXAMPLE DESIRED STATE CONFIGURATION (DSC)
# This is an example of how to deploy an Automation Account along with a basic DSC to support applying and
# enforcing a configuration on a basic virtual machine.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_automationaccount_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  version = "~>2.20"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 0.12"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_nsg_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "automation_account_dsc_rg" {
  name     = "${var.resource_group_name}-${var.postfix}"
  location = var.location
}

resource "random_password" "virtual_machine" {
  length           = 16
  override_special = "-_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}

resource "azurerm_automation_account" "automation_account" {
  name                = "${var.automation_account_name}-${var.postfix}"
  sku_name            = "Basic"
  location            = azurerm_resource_group.automation_account_dsc_rg.location
  resource_group_name = azurerm_resource_group.automation_account_dsc_rg.name
}

resource "azurerm_automation_dsc_configuration" "sample_dsc" {
  name                    = var.sample_dsc_name
  resource_group_name     = azurerm_resource_group.automation_account_dsc_rg.name
  automation_account_name = azurerm_automation_account.automation_account.name
  location                = azurerm_resource_group.automation_account_dsc_rg.location
  content_embedded        = file(var.sample_dsc_path)
  depends_on              = [azurerm_automation_account.automation_account]
}

# # ---------------------------------------------------------------------------------------------------------------------
# # COMPILE THE SAMPLE DSC IN THE AUTOMATION ACCOUNT
# # The Terraform `null_resource` signs in to Azure.  The second performs the compliation
# # Compilation is triggered on every one to ensure the latest changes are always applied
# # WARNING: ConvertTo-SecureString -String '${var.POWERSHELL_CLIENT_SECRET}, will expose the actual secret value in the logs.
# # ---------------------------------------------------------------------------------------------------------------------
resource "null_resource" "azure_signin_pwsh" {
  provisioner "local-exec" {
    command     = "$User = '${var.client_id}' ; $Pword =  ConvertTo-SecureString -String '${var.client_secret}' -AsPlainText -Force ; $Credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $User, $PWord ; Connect-AzAccount -Environment ${var.cloud_environment} -Credential $Credential -Tenant '${var.ARM_TENANT_ID}' -ServicePrincipal"
    interpreter = ["pwsh", "-Command"]
  }
  triggers = {
    always_run = timestamp()
  }
  depends_on = [azurerm_automation_dsc_configuration.sample_dsc]
}

resource "null_resource" "compileSampleDSC" {
  provisioner "local-exec" {
    command     = "Start-AzAutomationDscCompilationJob -ResourceGroupName  ${azurerm_resource_group.automation_account_dsc_rg.name} -AutomationAccountName ${azurerm_automation_account.automation_account.name} -ConfigurationName ${var.sample_dsc_name}"
    interpreter = ["pwsh", "-Command"]
  }
  triggers = {
    always_run = timestamp()
  }
  depends_on = [null_resource.azure_signin_pwsh]
}

resource "azurerm_automation_certificate" "automation_account_certificate" {
  name                    = "${var.automation_run_as_certificate_name}-${var.postfix}"
  resource_group_name     = azurerm_resource_group.automation_account_dsc_rg.name
  automation_account_name = azurerm_automation_account.automation_account.name

  description = var.automation_run_as_certificate_name
  # Certificate must be in .pfx format without a password encoded in base64
  base64 = var.automation_run_as_certificate_base64
}

resource "azurerm_automation_connection" "automation_account_connection" {
  name                    = "${var.automation_run_as_connection_name}-${var.postfix}"
  resource_group_name     = azurerm_resource_group.automation_account_dsc_rg.name
  automation_account_name = azurerm_automation_account.automation_account.name
  type                    = var.automation_run_as_connection_type

  values = {
    "ApplicationId" : var.client_id
    "TenantId" : var.ARM_TENANT_ID
    "SubscriptionId" : var.ARM_SUBSCRIPTION_ID
    "CertificateThumbprint" : var.automation_run_as_certificate_thumbprint
  }
}


# TEST VM RESOURCES
resource "azurerm_virtual_network" "vm_test" {
  name                = "sampledscvnet-${var.postfix}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.automation_account_dsc_rg.location
  resource_group_name = azurerm_resource_group.automation_account_dsc_rg.name
}

resource "azurerm_subnet" "vm_test" {
  name                 = "sampledscvmsubnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.automation_account_dsc_rg.name
  virtual_network_name = azurerm_virtual_network.vm_test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "vm_test" {
  name                = "sampledscni-${var.postfix}"
  location            = azurerm_resource_group.automation_account_dsc_rg.location
  resource_group_name = azurerm_resource_group.automation_account_dsc_rg.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.vm_test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "vm_test" {
  name                  = "${var.vm_host_name}-${var.postfix}"
  location              = azurerm_resource_group.automation_account_dsc_rg.location
  resource_group_name   = azurerm_resource_group.automation_account_dsc_rg.name
  network_interface_ids = [azurerm_network_interface.vm_test.id]
  vm_size               = "Standard_F2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "${var.vm_host_name}-${var.postfix}"
    admin_username = "testadmin"
    admin_password = random_password.virtual_machine.result
  }

  os_profile_windows_config {
    provision_vm_agent = true
  }
}

resource "azurerm_virtual_machine_extension" "vm_test" {
  name                 = var.sample_dsc_name
  virtual_machine_id   = azurerm_virtual_machine.vm_test.id
  publisher            = "Microsoft.Powershell"
  type                 = "DSC"
  type_handler_version = "2.77"
  depends_on           = [azurerm_virtual_machine.vm_test]

  settings = <<SETTINGS
    {
      "configurationArguments": {
          "RegistrationUrl": "${azurerm_automation_account.automation_account.dsc_server_endpoint}",
          "NodeConfigurationName": "${var.sample_dsc_configuration_name}"
      }
    }
  SETTINGS

  protected_settings = <<PROTECTED_SETTINGS
    {
      "configurationArguments": {
        "registrationKey": {
          "userName": "NOT_USED",
          "Password": "${azurerm_automation_account.automation_account.dsc_primary_access_key}"
        }
      }
    }
  PROTECTED_SETTINGS
}
