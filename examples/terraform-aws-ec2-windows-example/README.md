# Windows Instance Example

This folder provides a Packer template that can be used to build an Amazon Machine Image (AMI) of a Windows 2016 Server that comes pre-installed with: 

- The [Chocolately package manager](https://chocolatey.org/why-chocolatey) which makes it easy to install additional software packages onto Windows
- Git
- Python 3

In addition, this folder provides an example of how to launch a Windows instance based off this AMI that can be connected to via a Remote Desktop Protocol (RDP) client for the purposes of testing software or experimentation. 

This setup is ideal for "hot-reloading" code that you're actively developing and testing it against the Windows server. You can develop your code in your usual environment, perhaps a Mac or Linux laptop, yet see your changes reflected on the Windows server in seconds, by sharing a folder from your development machine with the Windows server via the RDP client. 

## Quick start 

Pre-requistes: 

- [Packer version v1.8.1 or newer](https://github.com/hashicorp/packer)
- [Terraform v1.0 or newer](https://github.com/hashicorp/terraform)
- An AWS account with valid security credentials

First, we'll build the AMI for the Windows Instance. Change into the packer directory: 

`cd packer` 

In order to build an Amazon Machine Image with Packer, you'll need to export your AWS account credentials. You can export your AWS credentials as the environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`. 

For more information on authenticating to your AWS account from the command line, see our blog post [Authenticating to AWS with Environment Variables](https://blog.gruntwork.io/authenticating-to-aws-with-environment-variables-e793d6f6d02e).

With your credentials properly exported, you can now run the packer build: 

`packer build build.pkr.hcl`

This may take upwards of 25 minutes to complete, but generally completes in about 5 minutes. Keep an eye on your EC2 dashboard and ensure that you have selected the correct region and that you are on the AMI view. Once your AMI status has changed from "Pending" to "Available", you can copy your AMI ID. 

Create a new file named `terraform.tfvars` in this same directory and enter the following variables: 

```hcl
ami_id           = "<the AMI ID you copied in the previous step>"
region           = "us-east-1"
root_volume_size = 100
```
Save the file. 

You're now ready to run terraform plan and check the output before proceeding: 

`terraform plan`

Take a look at the plan output and ensure everything looks correct. You should see a single EC2 instance being created along with supporting resources such as a security group and security group rules. 

Once you're satisfied that the plan looks good, run terraform apply to create the infrastructure: 

`terraform apply --auto-approve`

Once your resources apply successfully you'll see a similar output message containing the public IPv4 address of your Windows instance: 

`instance_ip = "35.84.139.82"`

