output "topic_name" {
  description = "The name of the Pub/Sub topic created."
  value       = google_pubsub_topic.example.name
}

output "subscription_name" {
  description = "The name of the Pub/Sub subscription created."
  value       = google_pubsub_subscription.example.name
}
