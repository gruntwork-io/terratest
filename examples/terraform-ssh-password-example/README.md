# Terraform SSH Password Example

This folder contains a simple Terraform module that deploys resources in [AWS](https://aws.amazon.com/) to demonstrate
how you can use Terratest to write automated tests for your AWS Terraform code. This module deploys an [EC2
Instance](https://aws.amazon.com/ec2/) with a public IP in the AWS region specified in the `aws_region` variable. The
EC2 Instance allows SSH requests on the port specified by the `ssh_port` variable, and is configured with a user data
script so that it will accept passwords for authentication.

Please note that the Terraform deployment outlined in
[the example directory for this test](/examples/terraform-ssh-password-example) will expect a default VPC to exist in
the target region for deployments to go into.

If this default VPC has been deleted from your AWS account, it can be recreated with the following command:

``` shell
$ aws ec2 create-default-vpc --region eu-west-2
{
    "Vpc": {
        "CidrBlock": "172.31.0.0/16",
        "InstanceTenancy": "default",
        "IsDefault": true,
        "State": "pending",
        ...
```

Check out [test/terraform_ssh_password_example_test.go](/test/terraform_ssh_password_example_test.go) to see how you
can write automated tests for this module.

Note that the example in this module is still fairly simplified, as the EC2 Instance doesn't do a whole lot! For a more
complicated, real-world, end-to-end example of a Terraform module and web server, see
[terraform-packer-example](/examples/terraform-packer-example).

**WARNING**: This module and the automated tests for it deploy real resources into your AWS account which can cost you
money. The resources are all part of the [AWS Free Tier](https://aws.amazon.com/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all AWS charges.

## Running this module manually

1. Sign up for [AWS](https://aws.amazon.com/).
1. Configure your AWS credentials using one of the [supported methods for AWS CLI
   tools](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html), such as setting the
   `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables. If you're using the `~/.aws/config` file for profiles then export `AWS_SDK_LOAD_CONFIG` as "True".
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Sign up for [AWS](https://aws.amazon.com/).
1. Configure your AWS credentials using one of the [supported methods for AWS CLI
   tools](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html), such as setting the
   `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables. If you're using the `~/.aws/config` file for profiles then export `AWS_SDK_LOAD_CONFIG` as "True".
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformSshPasswordExample`
