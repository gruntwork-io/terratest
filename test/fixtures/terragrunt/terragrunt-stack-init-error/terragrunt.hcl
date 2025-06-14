terraform {
  source = "..//terragrunt-stack-init-error"
  extra_arguments "common_vars" {
    commands = get_terraform_commands_that_need_vars()
    arguments = [
      "-var-file=terraform.tfvars"
    ]
  }
}

# This is intentionally invalid HCL syntax - missing closing brace
inputs = {
  test_var = "test_value"
  # Missing closing brace for the inputs block 