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
# DEPLOY A PUBSUB TOPIC AND SUBSCRIPTION
# See test/gcp/terraform_gcp_pubsub_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

# website::tag::1:: Deploy a Pub/Sub topic
resource "google_pubsub_topic" "example" {
  project = var.gcp_project_id
  name    = var.topic_name
}

# website::tag::2:: Create a Subscription to the topic so we can verify it
resource "google_pubsub_subscription" "example" {
  project = var.gcp_project_id
  name    = var.subscription_name
  topic   = google_pubsub_topic.example.name
}
