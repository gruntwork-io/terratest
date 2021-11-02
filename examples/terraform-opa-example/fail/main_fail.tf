module "instance_types" {
  # website::tag::1:: We expect this to fail the OPA check since it is sourcing the module locally and not from gruntwork-io GitHub.
  source     = "../pass"
  aws_region = var.aws_region
}
