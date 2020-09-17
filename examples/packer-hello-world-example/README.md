# Packer "Hello, World" Example

This folder contains the simplest possible Packer template—one that builds a Docker image with a text file that says
"Hello, World"!—to demonstrate how you can use Terratest to write automated tests for your Packer templates. 

Check out [test/packer_hello_world_example_test.go](/test/packer_hello_world_example_test.go) to see how you can write
automated tests for this simple template.




## Building the Packer template

1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Install [Docker](https://www.docker.com/) and make sure it's on your `PATH`.
1. Run `packer build build.json`.
1. Run `docker run -it --rm gruntwork/packer-hello-world-example cat /test.txt`.
1. You should see the text "Hello, World!"




## Running automated tests against the Packer template

1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Install [Docker](https://www.docker.com/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `go test -v -run TestPackerHelloWorldExample`
