# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------
provider "aws" {
  region     = var.region
  access_key = "dummy"
  secret_key = "dummy"

  # By default, terraform will interact with s3 in a way that uses subdomains
  # to specify information about the bucket. If you have issues with subdomains
  # of "localhost" not resolving properly, set this to true to tell terraform
  # to use url paths rather than subdomains for this information.
  # See: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#s3_use_path_style
  s3_use_path_style = "true"

  endpoints {
    sts = "http://localhost:5000"
    s3  = "http://localhost:5000"
  }
}

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A S3 BUCKET WITH VERSIONING ENABLED INCLUDING TAGS TO A LOCAL ENDPOINT
# See test/terraform_aws_endpoint_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

# Deploy and configure test S3 bucket with versioning and access log
resource "aws_s3_bucket" "test_bucket" {
  bucket = "${local.aws_account_id}-${var.tag_bucket_name}"

  tags = {
    Name        = var.tag_bucket_name
    Environment = var.tag_bucket_environment
  }
}

resource "aws_s3_bucket_logging" "test_bucket" {
  bucket        = aws_s3_bucket.test_bucket.id
  target_bucket = aws_s3_bucket.test_bucket_logs.id
  target_prefix = "TFStateLogs/"
}

resource "aws_s3_bucket_versioning" "test_bucket" {
  bucket = aws_s3_bucket.test_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_acl" "test_bucket" {
  bucket = aws_s3_bucket.test_bucket.id
  acl    = "private"
}


# Deploy S3 bucket to collect access logs for test bucket
resource "aws_s3_bucket" "test_bucket_logs" {
  bucket = "${local.aws_account_id}-${var.tag_bucket_name}-logs"

  tags = {
    Name        = "${local.aws_account_id}-${var.tag_bucket_name}-logs"
    Environment = var.tag_bucket_environment
  }

  force_destroy = true
}

resource "aws_s3_bucket_acl" "test_bucket_logs" {
  bucket = aws_s3_bucket.test_bucket_logs.id
  acl    = "log-delivery-write"
}

# Configure bucket access policies

resource "aws_s3_bucket_policy" "bucket_access_policy" {
  count  = var.with_policy ? 1 : 0
  bucket = aws_s3_bucket.test_bucket.id
  policy = data.aws_iam_policy_document.s3_bucket_policy.json
}

data "aws_iam_policy_document" "s3_bucket_policy" {
  statement {
    effect = "Allow"
    principals {
      # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
      # force an interpolation expression to be interpreted as a list by wrapping it
      # in an extra set of list brackets. That form was supported for compatibility in
      # v0.11, but is no longer supported in Terraform v0.12.
      #
      # If the expression in the following list itself returns a list, remove the
      # brackets to avoid interpretation as a list of lists. If the expression
      # returns a single list item then leave it as-is and remove this TODO comment.
      identifiers = [local.aws_account_id]
      type        = "AWS"
    }
    actions   = ["*"]
    resources = ["${aws_s3_bucket.test_bucket.arn}/*"]
  }

  statement {
    effect = "Deny"
    principals {
      identifiers = ["*"]
      type        = "AWS"
    }
    actions   = ["*"]
    resources = ["${aws_s3_bucket.test_bucket.arn}/*"]

    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values = [
        "false",
      ]
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# LOCALS
# Used to represent any data that requires complex expressions/interpolations
# ---------------------------------------------------------------------------------------------------------------------

data "aws_caller_identity" "current" {
}

locals {
  aws_account_id = data.aws_caller_identity.current.account_id
}

