# Kubernetes "Hello, World" Example

This folder contains the most minimal Kubernetes resource config—which deploys a simple webapp that responds with
"Hello, World!"—to demonstrate how you can use Terratest to write automated tests for Kubernetes.

Check out [test/kubernetes_hello_world_example_test.go](/test/kubernetes_hello_world_example_test.go) to see how you can 
write automated tests for this simple resource config.




## Deploying the Kubernetes resource

1. Setup a Kubernetes cluster. We recommend using a local version:
    - [Kubernetes on Docker For Mac](https://docs.docker.com/docker-for-mac/kubernetes/)
    - [Kubernetes on Docker For Windows](https://docs.docker.com/docker-for-windows/kubernetes/)
    - [minikube](https://github.com/kubernetes/minikube)
1. Install and setup [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to talk to the deployed
   Kubernetes cluster.
1. Run `kubectl apply -f hello-world-deployment.yml`




## Running automated tests against this Kubernetes deployment

1. Setup a Kubernetes cluster. We recommend using a local version:
    - [Kubernetes on Docker For Mac](https://docs.docker.com/docker-for-mac/kubernetes/)
    - [Kubernetes on Docker For Windows](https://docs.docker.com/docker-for-windows/kubernetes/)
    - [minikube](https://github.com/kubernetes/minikube)
1. Install and setup [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to talk to the deployed
   Kubernetes cluster.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `go test -v -tags kubernetes -run TestKubernetesHelloWorldExample`
