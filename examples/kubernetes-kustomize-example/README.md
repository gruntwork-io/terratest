# Kubernetes Kustomize Example

This folder contains a minimal Kubernetes resource config file to demonstrate how you can use Terratest to write
automated tests for Kubernetes.

This resource file deploys an nginx container as a single pod deployment with a node port service attached to it.

See the corresponding terratest code for an example of how to test this resource config:
- [kubernetes_kustomize_example_test.go](../../test/kubernetes_kustomize_example_test.go)


## Deploying the Kubernetes resource

1. Setup a Kubernetes cluster. We recommend using a local version:
    - [minikube](https://github.com/kubernetes/minikube)
    - [Kubernetes on Docker For Mac](https://docs.docker.com/docker-for-mac/kubernetes/)
    - [Kubernetes on Docker For Windows](https://docs.docker.com/docker-for-windows/kubernetes/)

1. Install and setup [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to talk to the deployed
   Kubernetes cluster.
1. Run `kubectl apply -k .`


## Running automated tests against this Kubernetes deployment

1. Setup a Kubernetes cluster. We recommend using a local version:
    - [minikube](https://github.com/kubernetes/minikube)
    - [Kubernetes on Docker For Mac](https://docs.docker.com/docker-for-mac/kubernetes/)
    - [Kubernetes on Docker For Windows](https://docs.docker.com/docker-for-windows/kubernetes/)

1. Install and setup [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to talk to the deployed
   Kubernetes cluster.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -tags kubernetes -run TestKubernetesKustomizeExample`
