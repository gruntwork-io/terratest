---
layout: collection-browser-doc
title: Quick start
category: getting-started
excerpt: Learn how to start with Terratest.
tags: ["quick-start"]
order: 101
nav_title: Documentation
nav_title_link: /docs/
custom_js:
  - examples
  - prism
---

## Requirements

Terratest uses the Go testing framework. To use Terratest, you need to install:

- [Go](https://golang.org/) (requires version >=1.21.1)

## Setting up your project

The easiest way to get started with Terratest is to copy one of the examples and its corresponding tests from this
repo. This quick start section uses a Terraform example, but check out the [Examples]({{site.baseurl}}/examples/) section for other
types of infrastructure code you can test (e.g., Packer, Kubernetes, etc).

1. Create an `examples` and `test` folder.

1. Copy the folder including all the files from the [basic terraform example](https://github.com/gruntwork-io/terratest/tree/master/examples/terraform-basic-example/) into the `examples` folder.

1. Copy the [basic terraform example test](https://github.com/gruntwork-io/terratest/blob/master/test/terraform_basic_example_test.go) into the `test` folder.

1. To configure dependencies, run:

    ```bash
    cd test
    go mod init "<MODULE_NAME>"
    go mod tidy
    ```

    Where `<MODULE_NAME>` is the name of your module, typically in the format
    `github.com/<YOUR_USERNAME>/<YOUR_REPO_NAME>`.

1. To run the tests:

    ```bash
    cd test
    go test -v -timeout 30m
    ```

    *(See [Timeouts and logging]({{ site.baseurl }}/docs/testing-best-practices/timeouts-and-logging/) for why the `-timeout` parameter is used.)*


## Terratest intro

The basic usage pattern for writing automated tests with Terratest is to:

1. Write tests using Go’s built-in [package testing](https://golang.org/pkg/testing/): you create a file ending in `_test.go` and run tests with the `go test` command. E.g., `go test my_test.go`.
1. Use Terratest to execute your _real_ IaC tools (e.g., Terraform, Packer, etc.) to deploy _real_ infrastructure (e.g., servers) in a _real_ environment (e.g., AWS).
1. Use the tools built into Terratest to validate that the infrastructure works correctly in that environment by making HTTP requests, API calls, SSH connections, etc.
1. Undeploy everything at the end of the test.

To make this sort of testing easier, Terratest provides a variety of helper functions and patterns for common infrastructure testing tasks, such as testing Terraform code, testing Packer templates, testing Docker images, executing commands on servers over SSH, making HTTP requests, working with AWS APIs, and so on.


## Example #1: Terraform "Hello, World"

Let's start with the simplest possible [Terraform](https://www.terraform.io/) code, which just outputs the text, 
"Hello, World" (if you’re new to Terraform, check out our [Comprehensive Guide to 
Terraform](https://blog.gruntwork.io/a-comprehensive-guide-to-terraform-b3d32832baca)): 
 
{% include examples/explorer.html example_id='terraform-hello-world' file_id='terraform_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}
 
How can you test this code to be confident it works correctly? Well, let’s think about how you would test it manually:

1. Run `terraform init` and `terraform apply` to execute the code.
1. When `apply` finishes, check that the output variable says, "Hello, World".
1. When you're done testing, run `terraform destroy` to clean everything up.
 
Using Terratest, you can write an automated test that performs the exact same steps! Here’s what the code looks like:
 
{% include examples/explorer.html example_id='terraform-hello-world' file_id='test_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}
 
This code does all the steps we mentioned above, including running `terraform init`, `terraform apply`, reading the 
output variable using `terraform output`, checking its value is what we expect, and running `terraform destroy` 
(using [`defer`](https://blog.golang.org/defer-panic-and-recover) to run it at the end of the test, whether the test 
succeeds or fails). If you put this code in a file called `terraform_hello_world_example_test.go`, you can run it by 
executing `go test`, and you’ll see output that looks like this (truncated for readability):

```
$ go test -v
=== RUN   TestTerraformHelloWorldExample
Running command terraform with args [init]
Initializing provider plugins...
[...]
Terraform has been successfully initialized!
[...]
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
Outputs:
hello_world = "Hello, World!"
[...]
Running command terraform with args [destroy -force -input=false]
[...]
Destroy complete! Resources: 2 destroyed.
--- PASS: TestTerraformHelloWorldExample (149.36s)
```

Success! 

## Example #2: Terraform and AWS

Let's now try out a more realistic Terraform example. Here is some Terraform code that deploys a simple web server in 
AWS:

{% include examples/explorer.html example_id='aws-hello-world' file_id='terraform_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}

The code above deploys an [EC2 Instance](https://aws.amazon.com/ec2/) that is running an Ubuntu 
[Amazon Machine Image (AMI)](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html). To keep this example 
simple, we specify a [User Data](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-api-cli) 
script that, while the server is booting, fires up a dirt-simple web server that returns “Hello, World” on port 8080.

How can you test this code to be confident it works correctly? Well, let’s again think about how you would test it 
manually:

1. Run `terraform init` and `terraform apply` to deploy the web server into your AWS account.
1. When `apply` finishes, get the IP of the web server by reading the `public_ip` output variable.
1. Open the IP in your web browser with port 8080 and make sure it says “Hello, World”. Note that it can take 1–2 
   minutes for the server to boot up, so you may have to retry a few times.
1. When you’re done testing, run `terraform destroy` to clean everything up.

Here's how we can automate the steps above using Terratest:

{% include examples/explorer.html example_id='aws-hello-world' file_id='test_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}

This test code runs `terraform init` and `terraform apply`, reads the server IP using `terraform output`, makes HTTP 
requests to the web server (including plenty of retries to account for the server taking time to boot), checks the HTTP
response is what we expect, and then runs `terraform destroy` at the end. If you put this code in a file called 
`terraform_aws_hello_world_example_test.go`, you can run just this test by passing the `-run` argument to `go test` as 
follows:

```
$ go test -v -run TestTerraformAwsHelloWorldExample -timeout 30m
=== RUN   TestTerraformAwsHelloWorldExample
Running command terraform with args [init]
Initializing provider plugins...
[...]
Terraform has been successfully initialized!
[...]
Running command terraform with args [apply -auto-approve]
aws_instance.example: Creating...
  associate_public_ip_address:       "" => "<computed>"
  availability_zone:                 "" => "<computed>"
  ephemeral_block_device.#:          "" => "<computed>"
  instance_type:                     "" => "t2.micro"
  key_name:                          "" => "<computed>"
[...]
Apply complete! Resources: 2 added, 0 changed, 0 destroyed.
Outputs:
public_ip = 52.67.41.31
[...]
Making an HTTP GET call to URL http://52.67.41.31:8080
dial tcp 52.67.41.31:8080: getsockopt: connection refused.
Sleeping for 5s and will try again.
Making an HTTP GET call to URL http://52.67.41.31:8080
dial tcp 52.67.41.31:8080: getsockopt: connection refused.
Sleeping for 5s and will try again.
Making an HTTP GET call to URL http://52.67.41.31:8080
Success!
[...]
Running command terraform with args [destroy -force -input=false]
[...]
Destroy complete! Resources: 2 destroyed.
--- PASS: TestTerraformAwsHelloWorldExample (149.36s)
```

Success! Now, every time you make a change to this Terraform code, the test code can run and make sure your web server 
works as expected.

Note that in the `go test` command above, we set `-timeout 30m`. This is because Go sets a default test time out of 10
minutes, and if your test take longer than that to run, Go will panic, and kill the test code part way through. This is
not only annoying, but also prevents the clean up code from running (the `terraform destroy`), leaving you with lots of
resources hanging in your AWS account. To prevent this, we always recommend setting a high test timeout; the test above
doesn't actually take anywhere near 30 minutes (typical runtime is ~3 minutes), but we give lots of extra buffer to be
extra sure that the test always has a chance to finish cleanly. 

## Example #3: Docker

You can use Terratest for testing a variety of infrastructure code, not just Terraform. For example, you can use it to
test your [Docker](https://www.docker.com/) images:

{% include examples/explorer.html example_id='docker-hello-world' file_id='docker_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}

The `Dockerfile` above creates a simple Docker image that uses Ubuntu 18.04 as a base and writes the text "Hello, World!" 
to a text file. At this point, you should already know the drill. First, let's think through how you'd test this 
`Dockerfile` manually:

1. Run `docker build` to build the Docker image.
1. Run the image via `docker run`.
1. Check that the running Docker container has a text file with the text "Hello, World!" in it.

Here's how you can use Terratest to automate this process:  

{% include examples/explorer.html example_id='docker-hello-world' file_id='test_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}

Instead of using Terraform helpers, this test code uses Terratest's Docker helpers to run `docker build`, `docker run`,
and check the contents of the text file. As before, you can run this test using `go test`!

## Example #4: Kubernetes

Terratest also provides helpers for testing your [Kubernetes](https://kubernetes.io/) code. For example, here's a 
Kubernetes manifest you might want to test:

{% include examples/explorer.html example_id='kubernetes-hello-world' file_id='k8s_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}

This manifest deploys the [Docker training webapp](https://hub.docker.com/r/training/webapp/), a simple app that 
responds with the text "Hello, World!", as a Kubernetes Deployment and exposes it to the outside world on port 5000 
using a `LoadBalancer`.

To test this code manually, you would:

1. Run `kubectl apply` to deploy the Docker training webapp.
1. Use the Kubernetes APIs to figure out the endpoint to hit for the load balancer.
1. Open the endpoint in your web browser on port 5000 and make sure it says “Hello, World”. Note that, depending on 
   your Kubernetes cluster, it could take a minute or two for the Docker container to come up, so you may have to retry 
   a few times.
1. When you're done testing, run `kubectl delete` to clean everything up.

Here's how you automate this process with Terratest:

{% include examples/explorer.html example_id='kubernetes-hello-world' file_id='test_code' class='wide quick-start-examples' skip_learn_more=true skip_view_on_github=true skip_tags=true %}

The test code above uses Kuberenetes helpers built into Terratest to run `kubectl apply`, wait for the service to come
up, get the service endpoint, make HTTP requests to the service (with plenty of retries), check the response is what
we expect, and runs `kubectl delete` at the end. You run this test with `go test` as well! 


## Give it a shot!

The above is just a small taste of what you can do with [Terratest](https://github.com/gruntwork-io/terratest). To 
learn more:

1. Check out the [examples]({{site.baseurl}}/examples/) and the corresponding automated tests for those examples for fully working (and tested!) sample code.
1. Browse through the list of [Terratest packages]({{site.baseurl}}/docs/getting-started/packages-overview/) to get a sense of all the tools available in Terratest.
1. Read our [Testing Best Practices Guide]({{site.baseurl}}/docs/#testing-best-practices).
1. Check out real-world examples of Terratest usage in our open source infrastructure modules: [Consul](https://github.com/hashicorp/terraform-aws-consul), [Vault](https://github.com/hashicorp/terraform-aws-vault), [Nomad](https://github.com/hashicorp/terraform-aws-nomad).

Happy testing!
