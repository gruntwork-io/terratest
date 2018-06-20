# Terratest

[![Go Report Card](https://goreportcard.com/badge/github.com/gruntwork-io/terratest)](https://goreportcard.com/report/github.com/gruntwork-io/terratest)

Terratest is a Go library that makes it easier to write automated tests for your infrastructure code. It provides a
variety of helper functions and patterns for common infrastructure testing tasks, including:

- Testing Terraform code
- Testing Packer templates
- Testing Docker images
- Executing commands on servers over SSH
- Working with AWS APIs
- Making HTTP requests
- Running shell commands
- And much more

Terratest was developed at [Gruntwork](https://gruntwork.io/) to help maintain the [Infrastructure as Code
Library](https://gruntwork.io/infrastructure-as-code-library/), which contains over 250,000 lines of code written
in Terraform, Go, Python, and Bash, and is used in production by hundreds of companies. See [Open sourcing Terratest: a
swiss army knife for testing infrastructure code](https://blog.gruntwork.io/open-sourcing-terratest-a-swiss-army-knife-for-testing-infrastructure-code-5d883336fcd5)
for more info.




## Introduction

The basic usage pattern for writing automated tests with Terratest is to:

1.  Write tests using Go's built-in [package testing](https://golang.org/pkg/testing/): you create a file ending in
    `_test.go` and run tests with the `go test` command.
1.  Use Terratest to execute your _real_ IaC tools (e.g., Terraform, Packer, etc.) to deploy _real_ infrastructure
    (e.g., servers) in a _real_ environment (e.g., AWS).
1.  Validate that the infrastructure works correctly in that environment by making HTTP requests, API calls, SSH
    connections, etc.
1.  Undeploy everything at the end of the test.

Here's a simple example of how to test some Terraform code:

```go
terraformOptions := &terraform.Options {
  // The path to where your Terraform code is located
  TerraformDir: "../examples/terraform-basic-example",
}

// At the end of the test, run `terraform destroy` to clean up any resources that were created
defer terraform.Destroy(t, terraformOptions)

// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
terraform.InitAndApply(t, terraformOptions)

// Validate your code works as expected
validateServerIsWorking(t, terraformOptions)
```




## Install

Prerequisite: install [Go](https://golang.org/).

To add Terratest to your projects, we recommend using a Go dependency manager such as
[dep](https://github.com/golang/dep) to add the packages you wish to use (see [package by package overview](#package-by-package-overview) for the list). For example, to add the `terraform` package:

```bash
dep ensure -add github.com/gruntwork-io/terratest/modules/terraform
```

Alternatively, you can use `go get`:

```bash
go get github.com/gruntwork-io/terratest/modules/terraform
```




## Examples

The best way to learn how to use Terratest is through examples.

First, check out the [examples folder](/examples) for different types of infrastructure code you may want to test,
such as:

1.  [Basic Terraform Example](/examples/terraform-basic-example): A simple "Hello, World" Terraform configuration.
1.  [HTTP Terraform Example](/examples/terraform-http-example): A more complicated Terraform configuration that deploys
    a simple web server that responds to HTTP requests in AWS.
1.  [Basic Packer Example](/examples/packer-basic-example): A simple Packer template for building an Amazon Machine
    Image (AMI).
1.  [Terraform Packer Example](/examples/terraform-packer-example): A more complicated example that shows how to use
    Packer to build an AMI with a web server installed and deploy that AMI in AWS using Terraform.

Next, head over to the [test folder](/test) to see how you can use Terraform to test each of these examples:

1.  [terraform_basic_example_test.go](/test/terraform_basic_example_test.go): Use Terratest to run `terraform apply` on
    the Basic Terraform Example and verify you get the expected outputs.
1.  [terraform_http_example_test.go](/test/terraform_http_example_test.go): Use Terratest to run `terraform apply` on
    the HTTP Terraform Example to deploy the web server, make HTTP requests to the web server to check that it is
    working correctly, and run `terraform destroy` to undeploy the web server.
1.  [packer_basic_example_test.go](/test/packer_basic_example_test.go): Use Terratest to run `packer build` to build an
    AMI and then use the AWS APIs to delete that AMI.
1.  [terraform_packer_example_test.go](/test/terraform_packer_example_test.go): Use Terratest to run `packer build` to
    build an AMI with a web server installed, deploy that AMI in AWS by running `terraform apply`, make HTTP requests to
    the web server to check that it is working correctly, and run `terraform destroy` to undeploy the web server.

Finally, to see some real-world examples of Terratest in action, check out some of our open source infrastructure
modules:

1.  [Consul](https://github.com/hashicorp/terraform-aws-consul)
1.  [Vault](https://github.com/hashicorp/terraform-aws-vault)
1.  [Nomad](https://github.com/hashicorp/terraform-aws-nomad)
1.  [Couchbase](https://github.com/gruntwork-io/terraform-aws-couchbase/)




## Package by package overview

Now that you've had a chance to browse the examples and their tests, here's an overview of the packages you'll find in
Terratest's [modules folder](/modules) and how they can help you test different types infrastructure:

| Package            | Description                                                                                                                                                                                                                                                                                          |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **aws**            | Functions that make it easier to work with the AWS APIs. Examples: find an EC2 Instance by tag, get the IPs of EC2 Instances in an ASG, create an EC2 KeyPair, look up a VPC ID.                                                                                                                     |
| **collections**    | Go doesn't have much of a collections library built-in, so this package has a few helper methods for working with lists and maps. Examples: subtract two lists from each other.                                                                                                                      |
| **docker**         | Functions that make it easier to work with Docker and Docker Compose. Examples: run `docker-compose` commands.                                                                                                                                                                                       |
| **files**          | Functions for manipulating files and folders. Examples: check if a file exists, copy a folder and all of its contents.                                                                                                                                                                               |
| **git**            | Functions for working with Git. Examples: get the name of the current Git branch.                                                                                                                                                                                                                    |
| **http-helper**    | Functions for making HTTP requests. Examples: make an HTTP request to a URL and check the status code and body contain the expected values, run a simple HTTP server locally.                                                                                                                        |
| **logger**         | A replacement for Go's `t.Log` and `t.Logf` that writes the logs to `stdout` immediately, rather than buffering them until the very end of the test. This makes debugging and iterating easier.                                                                                                      |
| **packer**         | Functions for working with Packer. Examples: run a Packer build and return the ID of the artifact that was created.                                                                                                                                                                                  |
| **random**         | Functions for generating random data. Examples: generate a unique ID that can be used to namespace resources so multiple tests running in parallel don't clash.                                                                                                                                      |
| **retry**          | Functions for retrying actions. Examples: retry a function up to a maximum number of retries, retry a function until a stop function is called, wait up to a certain timeout for a function to complete. These are especially useful when working with distributed systems and eventual consistency. |
| **shell**          | Functions to run shell commands. Examples: run a shell command and return its `stdout` and `stderr`.                                                                                                                                                                                                 |
| **ssh**            | Functions to SSH to servers. Examples: SSH to a server, execute a command, and return `stdout` and `stderr`.                                                                                                                                                                                         |
| **terraform**      | Functions for working with Terraform. Examples: run `terraform init`, `terraform apply`, `terraform destroy`.                                                                                                                                                                                        |
| **test_structure** | Functions for structuring your tests to speed up local iteration. Examples: break up your tests into stages so that any stage can be skipped by setting an environment variable.                                                                                                                     |




## Testing best practices

Testing infrastructure as code (IaC) is hard. With general purpose programming languages (e.g., Java, Python, Ruby),
you have a "localhost" environment where you can run and test the code before you commit. You can also isolate parts
of your code from external dependencies to create fast, reliable unit tests. With IaC, neither of these advantages is
typically available, as there isn't a "localhost" equivalent for most IaC code (e.g., I can't use Terraform to deploy
an AWS VPC on my own laptop) and there's no way to isolate your code from the outside world (i.e., the whole point of
a tool like Terraform is to make calls to AWS, so if you remove AWS, there's nothing left).

That means that most of the tests are going to be integration tests that deploy into a real AWS account. This makes
the tests effective at catching real-world bugs, but it also makes them much slower and more brittle. In this section,
we'll outline some best practices to minimize the downsides of this sort of testing.

1.  [Testing environment](#testing-environment)
1.  [Namespacing](#namespacing)
1.  [Cleanup](#cleanup)
1.  [Timeouts and logging](#timeouts-and-logging)
1.  [Avoid test caching](#avoid-test-caching)
1.  [Error handling](#error-handling)
1.  [Iterating locally using Docker](#iterating-locally-using-docker)
1.  [Iterating locally using test stages](#iterating-locally-using-test-stages)


### Testing environment

Since most automated tests written with Terratest can make potentially destructive changes in your environment, we
strongly recommend running tests in an environment that is totally separate from production. For example, if you are
testing infrastructure code for AWS, you should run your tests in a completely separate AWS account.

This means that you will have to write your infrastructure code in such a way that you can plug in ([dependency
injection](https://en.wikipedia.org/wiki/Dependency_injection) environment-specific details, such as account IDs,
domain names, IP addresses, etc. Adding support for this will typically make your code cleaner and more flexible.


### Namespacing

Just about all resources your tests create (e.g., servers, load balancers, machine images) should be "namespaced" with
a unique name to ensure that:

1.  You don't accidentally overwrite any "production" resources in that environment (though as mentioned in the previous
    section, your test environment should be completely isolated from prod anyway).
1.  You don't accidentally clash with other tests running in parallel.

For example, when deploying AWS infrastructure with Terraform, that typically means exposing variables that allow you
to configure auto scaling group names, security group names, IAM role names, and any other names that must be unique.

You can use Terratest's `random.UniqueId()` function to generate identifiers that are short enough to use in resource
names (just 6 characters) but random enough to make it unlikely that you'll have a conflict.

```go
uniqueId := random.UniqueId()
instanceName := fmt.Sprintf("terratest-http-example-%s", uniqueId)

terraformOptions := &terraform.Options {
  TerraformDir: "../examples/terraform-http-example",
  Vars: map[string]interface{} {
    "instance_name": instanceName,
  },
}

terraform.Apply(t, terraformOptions)
```


### Cleanup

Since automated tests with Terratest deploy real resources into real environments, you'll want to make sure your tests
always cleanup after themselves so you don't leave a bunch of resources lying around. Typically, you should use Go's
`defer` keyword to ensure that the cleanup code always runs, even if the test hits an error along the way.

For example, if your test runs `terraform apply`, you should run `terraform destroy` at the end to clean up:

```go
// Ensure cleanup always runs
defer terraform.Destroy(t, options)

// Deploy
terraform.Apply(t, options)

// Validate
checkServerWorks(t, options)
```

Of course, despite your best efforts, occasionally cleanup will fail, perhaps due to the CI server going down, or a bug
in your code, or a temporary network outage. To handle those cases, we run a tool called
[cloud-nuke](https://github.com/gruntwork-io/cloud-nuke) in our test AWS account on a nightly basis to clean up any
leftover resources.


### Timeouts and logging

Go's package testing has a default timeout of 10 minutes, after which it forcibly kills your tests (even your cleanup
code won't run!). It's not uncommon for infrastructure tests to take longer than 10 minutes, so you'll want to increase
this timeout:

```bash
go test -timeout 30m
```

Note that many CI systems will also kill your tests if they don't see any log output for a certain period of time
(e.g., 10 minutes in CircleCI). If you use Go's `t.Log` and `t.Logf` for logging in your tests, you'll find that these
functions buffer all log output until the very end of the test (see https://github.com/golang/go/issues/24929 for more
info). If you have a long-running test, this might mean you get no log output for more than 10 minutes, and the CI
system will shut down your tests. Moreover, if your test has a bug that causes it to hang, you won't see any log output
at all to help you debug it.

Therefore, we recommend instead using Terratest's `logger.Log` and `logger.Logf` functions, which log to `stdout`
immediately:

```go
func TestFoo(t *testing.T) {
  logger.Log(t, "This will show up in stdout immediately")
}
```

Finally, if you're testing multiple Go packages, be aware that Go will buffer log output—even that sent directly to
`stdout` by `logger.Log` and `logger.Logf`—until all the tests in the package are done. This leads to the same
difficulties with CI servers and debugging. The workaround is to tell Go to test each package sequentially using the
`-p 1` flag:

```bash
go test -timeout 30m -p 1 ./...
```


### Avoid test caching

Since Go 1.10, test results are automatically [cached](https://golang.org/doc/go1.10#test). This can lead to Go not
running your tests again if you haven't changed any of the Go code. Since you're probably mainly manipulating Terraform
files, you should consider turning the caching of test results off. This ensures that the tests are run every time
you run `go test` and the result is not just read from the cache.

To turn caching off, you can use the `GOCACHE` environment variable and set it to `off`:

```shell
$ GOCACHE=off go test -timeout 30m -p 1 ./...
```


### Error handling

Just about every method `foo` in Terratest comes in two versions: `foo` and `fooE` (e.g., `terraform.Apply` and
`terraform.ApplyE`).

- `foo`: The base method takes a `t *testing.T` as an argument. If the method hits any errors, it calls `t.Fatal` to
  fail the test.

- `fooE`: Methods that end with the capital letter `E` always return an `error` as the last argument and never call
  `t.Fatal` themselves. This allows you to decide how to handle errors.

You will use the base method name most of the time, as it allows you to keep your code more concise by avoiding
`if err != nil` checks all over the place:

```go
terraform.Init(t, terraformOptions)
terraform.Apply(t, terraformOptions)
url := terraform.Output(t, terraformOptions, "url")
```

In the code above, if `Init`, `Apply`, or `Output` hits an error, the method will call `t.Fatal` and fail the test
immediately, which is typically the behavior you want. However, if you are _expecting_ an error and don't want it to
cause a test failure, use the method name that ends with a capital `E`:

```go
if _, err := terraform.InitE(t, terraformOptions); err != nil {
  // Do something with err
}

if _, err := terraform.ApplyE(t, terraformOptions); err != nil {
  // Do something with err
}

url, err := terraform.OutputE(t, terraformOptions, "url")
if err != nil {
  // Do something with err
}
```

As you can see, the code above is more verbose, but gives you more flexibility with how to handle errors.


### Iterating locally using Docker

For most infrastructure code, your only option is to deploy into a real environment such as AWS. However, if you're
writing scripts (i.e., Bash, Python, or Go), you should be able to test them locally using Docker. Docker containers
typically build 10x faster and start 100x faster than real servers, so using Docker for testing can help you iterate
much faster.

Here are some techniques we use with Docker:

- If your script is used in a Packer template, add a [Docker
  builder](https://www.packer.io/docs/builders/docker.html) to the template so you can create a Docker image from the
  same code. See the [Packer Docker Example](https://github.com/gruntwork-io/terratest/tree/master/examples/packer-docker-example) for working sample code.

- We have prebuilt Docker images for major Linux distros that have many important dependencies (e.g., curl, vim,
  tar, sudo) already installed. See the [test-docker-images folder](/test-docker-images) for more details.

- Create a `docker-compose.yml` to make it easier to run your Docker image with all the ports, environment variables,
  and other settings it needs. See the [Packer Docker Example](https://github.com/gruntwork-io/terratest/tree/master/examples/packer-docker-example) for working sample code.

- With scripts in Docker, you can replace _some_ real-world dependencies with mocks! One way to do this is to create
  some "mock scripts" and to bind-mount them in `docker-compose.yml` in a way that replaces the real dependency. For
  example, if your script calls the `aws` CLI, you could create a mock script called `aws` that shows up earlier in the
  `PATH`. Using mocks allows you to test 100% locally, without external dependencies such as AWS.


### Iterating locally using test stages

Most automated tests written with Terratest consist of multiple "stages", such as:

1.  Build an AMI using Packer
1.  Deploy the AMI using Terraform
1.  Validate that the AMI works as expected
1.  Undeploy the AMI using Terraform

Often, while testing locally, you'll want to re-run some subset of these stages over and over again: for example, you
might want to repeatedly run the validation step while you work out the kinks. Having to run _all_ of these stages
each time you change a single line of code can be very slow.

This is where Terratest's `test_structure` package comes in handy: it allows you to explicitly break up your tests into
stages and to be able to disable any one of those stages simply by setting an environment variable. Check out the
[terraform_packer_example_test.go](/test/terraform_packer_example_test.go) for working sample code.




## Alternative testing tools

1.  [A list of infrastructure testing tools](#a-list-of-infrastructure-testing-tools)
1.  [How Terratest compares to other testing tools](#how-terratest-compares-to-other-testing-tools)


### A list of infrastructure testing tools

Below is a list of other infrastructure testing tools you may wish to use in addition to Terratest. Check out [How
Terratest compares to other testing tools](#how-terratest-compares-to-other-testing-tools) to understand the trade-offs.

1.  [kitchen-terraform](https://github.com/newcontext-oss/kitchen-terraform)
1.  [rspec-terraform](https://github.com/bsnape/rspec-terraform)
1.  [serverspec](https://serverspec.org/)
1.  [inspec](https://www.inspec.io/)
1.  [Goss](https://github.com/aelsabbahy/goss)
1.  [awspec](https://github.com/k1LoW/awspec)
1.  [Terraform's acceptance testing framework](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#writing-an-acceptance-test)
1.  [ruby_terraform](https://github.com/infrablocks/ruby_terraform)


### How Terratest compares to other testing tools

Most of the other infrastructure testing tools we've seen are focused on making it easy to check the properties of a
single server or resource. For example, the various `xxx-spec` tools offer a nice, concise language for connecting to
a server and checking if, say, `httpd` is installed and running. These tools are effectively verifying that individual
"properties" of your infrastructure meet a certain spec.

Terratest approaches the testing problem from a different angle. The question we're trying to answer is, "does the
infrastructure actually work?" Instead of checking individual server properties (e.g., is `httpd` installed and
running), we'll actually make HTTP requests to the server and check that we get the expected response; or we'll store
data in a database and make sure we can read it back out; or we'll try to deploy a new version of a Docker container
and make sure the orchestration tool can roll out the new container with no downtime.

Moreover, we use Terratest not only with individual servers, but to test entire systems. For example, the automated
tests for the [Vault module](https://github.com/hashicorp/terraform-aws-vault/tree/master/modules) do the following:

1.  Use Packer to build an AMI.
1.  Use Terraform to create self-signed TLS certificates.
1.  Use Terraform to deploy all the infrastructure: a Vault cluster (which runs the AMI from the previous step), Consul
    cluster, load balancers, security groups, S3 buckets, and so on.
1.  SSH to a Vault node to initialize the cluster.
1.  SSH to all the Vault nodes to unseal them.
1.  Use the Vault SDK to store data in Vault.
1.  Use the Vault SDK to make sure you can read the same data back out of Vault.
1.  Use Terraform to undeploy and clean up all the infrastructure.

The steps above are exactly what you would've done to test the Vault module manually. Terratest helps automate this
process. You can think of Terratest as a way to do end-to-end, acceptance or integration testing, whereas most other
tools are focused on unit or functional testing.


### Why Terratest?

Our experience with building the [Infrastructure as Code Library](https://gruntwork.io/infrastructure-as-code-library/)
is that the _only_ way to create reliable, maintainable infrastructure code is to have a thorough suite of real-world,
end-to-end acceptance tests. Without these sorts of tests, you simply cannot be confident that the infrastructure code
actually works.

This is especially important with modern DevOps, as all the tools are changing so quickly. Terratest has helped us
catch bugs not only in our own code, but also in AWS, Azure, Terraform, Packer, Kafka, Elasticsearch, CircleCI, and
so on. Moreover, by running tests nightly, we're able to catch backwards incompatible changes and
regressions in our dependencies (e.g., backwards incompatibilities in new versions of Terraform) as early as possible.




## Developing Terratest

1.  [Contributing](#contributing)
1.  [Running tests](#running-tests)
1.  [Versioning](#versioning)


### Contributing

Contributions are very welcome! Check out the
[Contribution Guidelines](CONTRIBUTING.md) for instructions.


### Running tests

Terratest itself includes a number of automated tests.

**Note #1**: Some of these tests create real resources in an AWS account. That means they cost money to run, especially
if you don't clean up after yourself. Please be considerate of the resources you create and take extra care to clean
everything up when you're done!

**Note #2**: In order to run tests that access your AWS account, you will need to configure your [AWS CLI
credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html). For example, you could
set the credentials as the environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.

**Note #3**: Never hit `CTRL + C` or cancel a build once tests are running or the cleanup tasks won't run!

**Prerequisite**: Most the tests expect Terraform, Packer, and/or Docker to already be installed and in your `PATH`.

To run all the tests:

```bash
go test -v -timeout 30m -p 1 ./...
```

To run the tests in a specific folder:

```bash
cd "<FOLDER_PATH>"
go test -timeout 30m
```

To run a specific test in a specific folder:

```bash
cd "<FOLDER_PATH>"
go test -timeout30m -run "<TEST_NAME>"
```


### Versioning

This repo follows the principles of [Semantic Versioning](http://semver.org/). You can find each new release,
along with the changelog, in the [Releases Page](https://github.com/gruntwork-io/terratest/releases).

During initial development, the major version will be 0 (e.g., `0.x.y`), which indicates the code does not yet have a
stable API. Once we hit `1.0.0`, we will make every effort to maintain a backwards compatible API and use the MAJOR,
MINOR, and PATCH versions on each release to indicate any incompatibilities.




## License

This code is released under the Apache 2.0 License. Please see [LICENSE](LICENSE) and [NOTICE](NOTICE) for more details.

Copyright &copy; 2018 Gruntwork, Inc.
