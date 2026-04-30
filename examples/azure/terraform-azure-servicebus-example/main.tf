# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE SERVICE BUS
# This is an example of how to deploy an Azure service bus.
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

provider "azurerm" {
  features {}
}
# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "servicebus_rg" {
  name     = "terratest-sb-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# Define locals variables
# ---------------------------------------------------------------------------------------------------------------------
locals {
  topic_authorization_rules = flatten([
    for topic in var.topics : [
      for rule in topic.authorization_rules :
      merge(
        rule, {
          topic_name = topic.name
      })
    ]
  ])

  topic_subscriptions = flatten([
    for topic in var.topics : [
      for subscription in topic.subscriptions :
      merge(
        subscription, {
          topic_name = topic.name
      })
    ]
  ])

  topic_subscription_rules = flatten([
    for subscription in local.topic_subscriptions :
    merge({
      filter_type = ""
      sql_filter  = ""
      action      = ""
      }, subscription, {
      topic_name        = subscription.topic_name
      subscription_name = subscription.name
    })
    if subscription.filter_type != null
  ])
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE Service Bus Namespace
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_servicebus_namespace" "servicebus" {
  name                = "terratest-namespace-${var.namespace_name}"
  location            = azurerm_resource_group.servicebus_rg.location
  resource_group_name = azurerm_resource_group.servicebus_rg.name
  sku                 = var.sku
  tags                = var.tags
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE Service Bus Namespace Authorization Rule
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_servicebus_namespace_authorization_rule" "sbnamespaceauth" {
  count = length(var.namespace_authorization_rules)

  name         = var.namespace_authorization_rules[count.index].policy_name
  namespace_id = azurerm_servicebus_namespace.servicebus.id

  listen = var.namespace_authorization_rules[count.index].claims.listen
  send   = var.namespace_authorization_rules[count.index].claims.send
  manage = var.namespace_authorization_rules[count.index].claims.manage
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE Service Bus Topic
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_servicebus_topic" "sptopic" {
  count = length(var.topics)

  name         = var.topics[count.index].name
  namespace_id = azurerm_servicebus_namespace.servicebus.id

  requires_duplicate_detection = var.topics[count.index].requires_duplicate_detection
  default_message_ttl          = var.topics[count.index].default_message_ttl
  partitioning_enabled         = var.topics[count.index].enable_partitioning
  support_ordering             = var.topics[count.index].support_ordering
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE Service Bus Topic Authorization Rule
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_servicebus_topic_authorization_rule" "topicaauth" {
  count = length(local.topic_authorization_rules)

  name = local.topic_authorization_rules[count.index].policy_name
  topic_id = [for t in azurerm_servicebus_topic.sptopic :
  t.id if t.name == local.topic_authorization_rules[count.index].topic_name][0]

  listen = local.topic_authorization_rules[count.index].claims.listen
  send   = local.topic_authorization_rules[count.index].claims.send
  manage = local.topic_authorization_rules[count.index].claims.manage

  depends_on = [azurerm_servicebus_topic.sptopic]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE Service Bus Subscription
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_servicebus_subscription" "subscription" {
  count = length(local.topic_subscriptions)

  name = local.topic_subscriptions[count.index].name
  topic_id = [for t in azurerm_servicebus_topic.sptopic :
  t.id if t.name == local.topic_subscriptions[count.index].topic_name][0]

  max_delivery_count                   = local.topic_subscriptions[count.index].max_delivery_count
  lock_duration                        = local.topic_subscriptions[count.index].lock_duration
  forward_to                           = local.topic_subscriptions[count.index].forward_to
  dead_lettering_on_message_expiration = local.topic_subscriptions[count.index].dead_lettering_on_message_expiration

  depends_on = [azurerm_servicebus_topic.sptopic]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE Service Bus Subscription Rules
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_servicebus_subscription_rule" "subrules" {
  count = length(local.topic_subscription_rules)

  name = local.topic_subscription_rules[count.index].name
  subscription_id = [for s in azurerm_servicebus_subscription.subscription :
  s.id if s.name == local.topic_subscription_rules[count.index].subscription_name][0]
  filter_type = local.topic_subscription_rules[count.index].filter_type != "" ? "SqlFilter" : null
  sql_filter  = local.topic_subscription_rules[count.index].sql_filter
  action      = local.topic_subscription_rules[count.index].action

  depends_on = [azurerm_servicebus_subscription.subscription]
}
