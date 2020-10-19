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

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "instance_name" {
  description = "The Name to use for the Cloud Instance."
  type        = string
  default     = "terratest-example"
}

variable "machine_type" {
  description = "The Machine Type to use for the Cloud Instance."
  type        = string
  default     = "f1-micro"
}

variable "zone" {
  description = "The Zone to launch the Cloud Instance into."
  type        = string
  default     = "us-central1-a"
}

variable "bucket_name" {
  description = "The Name of the example Bucket to create."
  type        = string
  default     = "gruntwork-terratest-bucket"
}

variable "bucket_location" {
  description = "The location to store the Bucket. This value can be regional or multi-regional."
  type        = string
  default     = "US"
}

