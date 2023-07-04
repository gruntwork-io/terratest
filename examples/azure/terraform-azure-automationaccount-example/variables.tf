# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# Note that the TF_VAR_AUTOMATION_ACCOUNT_CLIENT_ID and TF_VAR_AUTOMATION_ACCOUNT_CLIENT_PASSWORD 
# variables represent a service principal in AAD
# In addition to created a service prinicpal for the automation run as connection, a certificate 
# must be uploaded to the service principal as a secret.
# The same certificate must also be uplaoded into the Azure Automation Account via Terraform 
# ---------------------------------------------------------------------------------------------------------------------

# ARM_CLIENT_ID
# ARM_CLIENT_SECRET
# ARM_SUBSCRIPTION_ID
# ARM_TENANT_ID
# ARM_ENVIRONMENT
# AZURE_ENVIRONMENT
# TF_VAR_cloud_environment
# TF_VAR_location
# TF_VAR_ARM_SUBSCRIPTION_ID
# TF_VAR_ARM_TENANT_ID
# TF_VAR_automation_run_as_certificate_thumbprint
# TF_VAR_client_id
# TF_VAR_client_secret
# TF_VAR_automation_run_as_certificate_thumbprint
# TF_VAR_automation_run_as_certificate_base64
# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "ARM_SUBSCRIPTION_ID" {
  type        = string
  description = "The Azure Subscription ID where the infrastructure will be deployed."
}
variable "ARM_TENANT_ID" {
  type        = string
  description = "The Azure Active Directory Tenant that the Azure Subscription is associated with."
}

variable "client_id" {
  sensitive   = true
  type        = string
  description = "The Azure Active Directory App registration client_id with access to the Azure Subscription."
}

variable "client_secret" {
  sensitive   = true
  type        = string
  description = "The secret for the Azure Active Directory App registration client_idwith access to the Azure Subscription."
}

variable "automation_run_as_certificate_base64" {
  type        = string
  description = "The certificate required for Automation Account 'Run As' account creation."
}

variable "automation_run_as_certificate_thumbprint" {
  type        = string
  description = "The thumbprint for the certificate required for Automation Account 'Run As' account creation."
}

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
  default     = "terratest-automationaccount-rg"
}

variable "location" {
  description = "The Azure region in which to deploy this sample"
  type        = string
  default     = "East US"
}

variable "cloud_environment" {
  description = "The Azure cloud where the command is executed"
  type        = string
  default     = "AzureCloud"
}

variable "automation_account_name" {
  description = "The name of the automation account that will be created in the resource group"
  type        = string
  default     = "terratest-AutomationAccount"
}

variable "automation_run_as_connection_name" {
  description = "The name of the automation run as connection that will be created in the resource group"
  type        = string
  default     = "terratest-AutomationRunAsConnectionName"
}

variable "automation_run_as_certificate_name" {
  description = "The name of the automation account run as connection certificate name"
  type        = string
  default     = "terratest-AutomationConnectionCertificateName"
}

variable "automation_run_as_certificate_path" {
  description = "The path to the automation Run As certificate .pfx file"
  type        = string
  default     = "./certificate/SPRunAsCert.pfx"
}

variable "automation_run_as_connection_type" {
  description = "The name of the automation account run as connection type"
  type        = string
  default     = "AzureServicePrincipal"
}

variable "sample_dsc_name" {
  description = "The name of the sample DSC configuration that contains the configuration that can be applied"
  type        = string
  default     = "SampleDSC"
}

variable "sample_dsc_path" {
  description = "The path to the sample dsc file in the repo"
  type        = string
  default     = "./dsc/SampleDSC.ps1"
}

variable "sample_dsc_configuration_name" {
  description = "The name of the DSC configuration to apply to the VM"
  type        = string
  default     = "SampleDSC.NotWebServer"
}

variable "vm_name" {
  description = "The name of the test VM where the DSC configuration will be applied"
  type        = string
  default     = "vm01"
}

variable "vm_host_name" {
  description = "The host name of the VM machine where the DSC configuration will be applied"
  type        = string
  default     = "dscnode"
}
