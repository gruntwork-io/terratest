terraform {
  required_providers {
    validation = {
      source  = "tlkamp/validation"
      version = "1.1.1"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.2.2"
    }
  }
}

# this data source will produce warning when `condition` is evaluated to `true`
data "validation_warning" "warn" {
  for_each  = toset([for i in range(10) : format("%02d", i)])
  condition = true
  summary   = "lorem ipsum ${each.value}"
  details   = "lorem ipsum dolor sit amet"
}

resource "null_resource" "empty" {}
