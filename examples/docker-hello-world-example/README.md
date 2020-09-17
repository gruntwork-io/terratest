# Docker "Hello, World" Example

This folder contains a `Dockerfile` to build a very simple Docker image—one which contains a text file with the
text "Hello, World"!—to demonstrate how you can use Terratest to write automated tests for your Docker images. 

Check out [test/docker_hello_world_example_test.go](/test/docker_hello_world_example_test.go) to see how you can write
automated tests for this simple Docker image.




## Building the Docker container

1. Install [Docker](https://www.docker.com/) and make sure it's on your `PATH`.
1. Run `docker build -t gruntwork/docker-hello-world-example .`.
1. Run `docker run -it --rm gruntwork/docker-hello-world-example cat /test.txt`.
1. You should see the text "Hello, World!"




## Running automated tests against the Docker container

1. Install [Docker](https://www.docker.com/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `go test -v -run TestDockerHelloWorldExample`
