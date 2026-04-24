# Terraform ASG SCP Example

This folder contains a simple Terraform module that deploys resources in [AWS](https://aws.amazon.com/) to demonstrate
how you can use Terratest to write automated tests for your AWS Terraform code. This module deploys an ASG with one instance.
The EC2 Instance allows SSH requests on the port specified by the `ssh_port` variable. 

Check out [test/terraform_scp_example_test.go](/test/terraform_scp_example_test.go) to see how you can write
automated tests for this module.

Note that the example in this module is still fairly simplified, as the EC2 Instance doesn't do a whole lot! For a more
complicated, real-world, end-to-end example of a Terraform module and web server, see
[terraform-packer-example](/examples/terraform-packer-example).

**WARNING**: This module and the automated tests for it deploy real resources into your AWS account which can cost you
money. The resources are all part of the [AWS Free Tier](https://aws.amazon.com/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all AWS charges.

## Overview 

When a test fails, it is often important to be able to quickly get to logs and config files from your deployed apps and services. Currently, getting at this information is a bit of a pain. Often times, it would be necessary to "catch it in the act". Usually this would require running tests and then "pausing"/not tearing down the infrastructure, ssh-ing to individual instances and then viewing the logs/config files that way. This in not very convenient and gets even more tricky when trying to get the same results for tests being executed by your CI server.

You can use terratest to help with this task by specifying `RemoteFileSpecification` structs that describe which files you want to copy from your instances:

```go
ctx := t.Context()

logstashSpec := aws.RemoteFileSpecification{
	SshUser:sshUserName,
	UseSudo:true,
	KeyPair:keyPair,
	LocalDestinationDir:filepath.Join("/tmp", "logs", t.Name(), "logstash"),
	AsgNames: strings.Split(strings.Replace(terraform.OutputRequiredContext(t, ctx, terraformOptions, "logstash_server_asg_names"), "\n", "", -1), ","),
	RemotePathToFileFilter: map[string][]string {
		"/var/log/logstash":{"*"},
		"/etc/logstash/conf.d" : {"*"},
	},
}
```

Once you've described what files you want, grabbing them from ASGs is simple with:
```go
aws.FetchFilesFromAsgsPContextE(t, ctx, awsRegion, &logstashSpec)
```

or directly from EC2 instances with:
```go
aws.FetchFilesFromInstanceContext(t, ctx, awsRegion, sshUserName, keyPair, appServerInstanceId, true, appServerConfig, filepath.Join("/tmp", "logs", t.Name(), "app_server"), []string{"*.yml", "caFile", "*.key", "*.pem"})
```

Finally, to put all of this together, in your go test you could do something like:

```go
defer test_structure.RunTestStage(t, "grab_logs", func() {
	if t.Failed() {
		takeElkMultiClusterLogSnapshot(t, examplesDir, awsRegion, "ubuntu")
	}
})
```

The code above will run at the very end of your test and grab a snapshot of all of the log descriptors you've defined specifically when your tests fail. We usually like to defer the logging snapshot right before we defer the terraform teardown. This way, if our tests fail, we are able to grab a snapshot of all the relevant logs and config files across our whole deployment!

You can even take this a step further and pair this with your CI's artifact storage mechanism to have all of your logs and config files attached to your broken CI build when tests failed. For example, with CircleCI, you could do something like:

```yml
      - store_artifacts:
          path: /tmp/logs
```

Now, to get at your logs when your tests fail, you will just be able to click the links right in your CI.

![logs](https://user-images.githubusercontent.com/34349331/46639252-086e0a00-cb33-11e8-8dd2-9be73ca2af56.gif)

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
1. `go test -v -run TestTerraformScpExample`
