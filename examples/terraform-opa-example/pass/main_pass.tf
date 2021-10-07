provider "aws" {
  region = var.aws_region
}

module "instance_types" {
  # website::tag::1:: We expect this to pass the OPA check since it is sourcing the module from gruntwork-io GitHub.
  source         = "git::git@github.com:gruntwork-io/terraform-aws-utilities.git//modules/instance-type?ref=v0.6.0"
  instance_types = ["t2.micro", "t3.micro"]
}
