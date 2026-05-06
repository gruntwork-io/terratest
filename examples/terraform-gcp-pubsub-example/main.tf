terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
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
