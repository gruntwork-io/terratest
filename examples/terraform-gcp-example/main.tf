# ---------------------------------------------------------------------------------------------------------------------
# CREATE A GOOGLE STORAGE BUCKET
# ---------------------------------------------------------------------------------------------------------------------

resource "google_storage_bucket" "test_bucket" {
    name     = "${var.bucket_name}"
    location = "${var.location}"
    project  = "${var.project}"
}
