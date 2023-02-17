# Helm Dependency Example

This folder contains a minimal helm chart to demonstrate how you can use Terratest to test your helm charts with dependencies.

There are two kinds of tests you can perform on a helm chart:

- Helm Template tests are tests designed to test the logic of the templates. These tests should run `helm template` with
  various input values and parse the yaml to validate any logic embedded in the templates (e.g by reading them in using
  client-go). Since templates are not statically typed, the goal of these tests is to promote fast cycle time

The helm chart deploys a single replica `Deployment` resource given the container image spec and a `Service` that
exposes it. This chart requires the `containerImageRepo` and `containerImageTag` input values.

See the corresponding terratest code for an example of how to test this chart:

- [helm_basic_example_template_test.go](/test/helm_basic_example_template_test.go): the template tests for this chart.

## Running automated tests against this Helm Chart

1. Install and setup [helm](https://docs.helm.sh/using_helm/#installing-helm)
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -tags helm -run TestHelmDependencyExampleTemplate` for the template test

**NOTE**: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly, helm
can overload the minikube system and thus interfere with the other kubernetes tests. Specifically, many of the tests
start to fail with `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes
tests and helm tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.
We recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.
