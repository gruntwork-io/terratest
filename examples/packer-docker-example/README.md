# Packer Docker Example

This folder contains a Packer template to demonstrate how you can use Terratest to write automated tests for your
Packer templates. The template creates an Ubuntu AMI with a simple web app (built on top of Ruby / Sinatra) installed.
This template _also_ creates a Docker image with the same web app installed, and contains a `docker-compose.yml` file
for running that Docker image. These allow you to test your Packer template completely locally, without having to
deploy to AWS.

Check out [test/packer_docker_example_test.go](/test/packer_docker_example_test.go) to see how you can write
automated tests for this simple template.

The Docker-based tests in this folder are in some sense "unit tests" for the Packer template. To see an example of
"integration tests" that deploy the AMI to AWS, check out the
[terraform-packer-example](/examples/terraform-packer-example).


## Installation steps
1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Install [Docker](https://www.docker.com/) and make sure it's on your `PATH`.


## Building a Docker image for local testing (Packer >= 1.7.0)
1. Run `packer init build.pkr.hcl`.
1. Run `packer build build.pkr.hcl`.


## Building a Docker image for local testing (Packer < 1.7.0)
1. Run `packer build build.json`.


## Run the container
1. Run `docker compose up`.
1. You should now be able to access the sample web app at http://localhost:8080




## Building an AMI for testing in AWS

1. Sign up for [AWS](https://aws.amazon.com/).
1. Configure your AWS credentials using one of the [supported methods for AWS CLI
   tools](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html), such as setting the
   `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Run `packer build -only=ubuntu-ami build.json`.




## Running automated tests locally against this Packer template

1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Install [Docker](https://www.docker.com/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestPackerDockerExampleLocal`




## Running automated tests in AWS against this Packer template

See [terraform-packer-example](/examples/terraform-packer-example).
