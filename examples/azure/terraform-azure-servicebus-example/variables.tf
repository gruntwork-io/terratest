# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# ARM_CLIENT_ID
# ARM_CLIENT_SECRET
# ARM_SUBSCRIPTION_ID
# ARM_TENANT_ID

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "location" {
  description = "The supported azure location where the resource exists"
  type        = string
  default     = "West US2"
}

variable "postfix" {
  description = "string mitigate resource name collisions."
  type        = string
  default     = "servicebus"
}

variable "namespace_name" {
  description = "The name of the namespace."
  type        = string
  default     = "testservicebus101"
}

variable "sku" {
  description = "The SKU of the namespace. The options are: `Basic`, `Standard`, `Premium`."
  type        = string
  default     = "Standard"
}

variable "tags" {
  description = " A mapping of tags to assign to the resource."
  type        = map(string)
  default     = {}
}

variable "namespace_authorization_rules" {
  description = "List of namespace authorization rules."
  type = list(object({
    policy_name = string
    claims      = object({ listen = bool, manage = bool, send = bool })
  }))
  default = []
}

variable "topics" {
  description = "topics list"
  type = list(object({
    name                         = string
    default_message_ttl          = string //ISO 8601 format
    enable_partitioning          = bool
    requires_duplicate_detection = bool
    support_ordering             = bool
    authorization_rules = list(object({
      policy_name = string
      claims      = object({ listen = bool, manage = bool, send = bool })

    }))
    subscriptions = list(object({
      name                                 = string
      max_delivery_count                   = number
      lock_duration                        = string //ISO 8601 format
      forward_to                           = string //set with the topic name that will be used for forwarding. Otherwise, set to ""
      dead_lettering_on_message_expiration = bool
      filter_type                          = string // SqlFilter is the only supported type now.
      sql_filter                           = string //Required when filter_type is set to SqlFilter
      action                               = string
    }))
  }))
  default = []
}
