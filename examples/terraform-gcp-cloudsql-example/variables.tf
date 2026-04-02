# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# You must define the following environment variables.
# ---------------------------------------------------------------------------------------------------------------------

# GOOGLE_CREDENTIALS
# or
# GOOGLE_APPLICATION_CREDENTIALS

variable "gcp_project_id" {
  description = "The ID of the GCP project in which these resources will be created."
}

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------
# (none)

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "instance_name" {
  description = "The name of the Cloud SQL instance to create."
  type        = string
  default     = "terratest-example-cloudsql"
}

variable "database_version" {
  description = "The database engine version to use (e.g. MYSQL_8_0, POSTGRES_14, SQLSERVER_2019_STANDARD)."
  type        = string
  default     = "POSTGRES_14"
}

variable "region" {
  description = "The GCP region in which to create the Cloud SQL instance."
  type        = string
  default     = "us-central1"
}

variable "tier" {
  description = "The machine type to use for the Cloud SQL instance."
  type        = string
  default     = "db-f1-micro"
}
