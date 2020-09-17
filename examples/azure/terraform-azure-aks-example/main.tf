# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE AKS CLUSTER
# This is an example of how to deploy an Azure AKS cluster with load balancer in front of the service 
# to handle providing the public interface into the cluster.
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ------------------------------------------------------------------------------

provider "azurerm" {
  version = "1.40.0"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "k8s" {
  name     = var.resource_group_name
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE KUBERNETES CLUSTER
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_kubernetes_cluster" "k8s" {
  name                = var.cluster_name
  location            = azurerm_resource_group.k8s.location
  resource_group_name = azurerm_resource_group.k8s.name
  dns_prefix          = var.dns_prefix

  linux_profile {
    admin_username = "ubuntu"

    ssh_key {
      key_data = file(var.ssh_public_key)
    }
  }

  default_node_pool {
    name       = "agentpool"
    node_count = var.agent_count
    vm_size    = "Standard_DS1_v2"
  }

  service_principal {
    client_id     = var.client_id
    client_secret = var.client_secret
  }

  tags = {
    Environment = "Development"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE KUBECONFIG FILE
# ---------------------------------------------------------------------------------------------------------------------

resource "local_file" "kubeconfig" {
  content  = azurerm_kubernetes_cluster.k8s.kube_config_raw
  filename = "kubeconfig"

  depends_on = [
    azurerm_kubernetes_cluster.k8s
  ]
}
