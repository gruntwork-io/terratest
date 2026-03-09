output "instance_name" {
  description = "The name of the Cloud SQL instance created."
  value       = google_sql_database_instance.example.name
}

output "database_version" {
  description = "The database engine version of the Cloud SQL instance."
  value       = google_sql_database_instance.example.database_version
}
