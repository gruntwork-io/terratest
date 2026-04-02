# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A CLOUD SQL INSTANCE
# See test/gcp/terraform_gcp_cloudsql_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

# website::tag::1:: Deploy a Cloud SQL instance
resource "google_sql_database_instance" "example" {
  project          = var.gcp_project_id
  name             = var.instance_name
  database_version = var.database_version
  region           = var.region

  settings {
    tier = var.tier
  }

  deletion_protection = false
}
