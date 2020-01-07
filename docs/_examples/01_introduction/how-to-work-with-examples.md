---
layout: collection-browser-doc
title: How to work with examples
category: introduction
excerpt: >-
  Terratest examples is the best way to start testing Terraform, Docker, Packer, Kubernetes, AWS, GCP, and more.
tags: ["introduction"]
order: 100
nav_title: Examples
nav_title_link: /examples/
---

The best way to learn how to use Terratest is through examples.

First, check out the [examples folder](/examples) for different types of infrastructure code you may want to test,
such as:

1.  [Basic Terraform Example](/examples/terraform-basic-example): A simple "Hello, World" Terraform configuration.
1.  [HTTP Terraform Example](/examples/terraform-http-example): A more complicated Terraform configuration that deploys
    a simple web server that responds to HTTP requests in AWS.
1.  [Basic Packer Example](/examples/packer-basic-example): A simple Packer template for building an Amazon Machine
    Image (AMI) or Google Cloud Platform Compute Image.
1.  [Terraform Packer Example](/examples/terraform-packer-example): A more complicated example that shows how to use
    Packer to build an AMI with a web server installed and deploy that AMI in AWS using Terraform.
1.  [Terraform GCP Example](/examples/terraform-gcp-example): A simple Terraform configuration that creates a GCP Compute Instance and Storage Bucket.
1.  [Terraform remote-exec Example](/examples/terraform-remote-exec-example): A terraform configuration that creates and
    AWS instance and then uses `remote-exec` to provision it.
1.  [Basic Kubernetes Example](/examples/kubernetes-basic-example): A minimal Kubernetes resource that deploys an
    addressable nginx instance.
1.  [Kubernetes RBAC Example](/examples/kubernetes-rbac-example): A Kubernetes resource config that creates a Namespace
    with a ServiceAccount that has admin permissions within the Namespace, but not outside.
1.  [Basic Helm Chart Example](/examples/helm-basic-example): A minimal helm chart that deploys a `Deployment` resource
    for the provided container image.

Next, head over to the [test folder](/test) to see how you can use Terratest to test each of these examples:

1.  [terraform_basic_example_test.go](/test/terraform_basic_example_test.go): Use Terratest to run `terraform apply` on
    the Basic Terraform Example and verify you get the expected outputs.
1.  [terraform_http_example_test.go](/test/terraform_http_example_test.go): Use Terratest to run `terraform apply` on
    the HTTP Terraform Example to deploy the web server, make HTTP requests to the web server to check that it is
    working correctly, and run `terraform destroy` to undeploy the web server.
1.  [packer_basic_example_test.go](/test/packer_basic_example_test.go): Use Terratest to run `packer build` to build an
    AMI and then use the AWS APIs to delete that AMI.
1.  [packer_gcp_basic_example_test.go](/test/packer_gcp_basic_example_test.go): Use Terratest to run `packer build`
    to build a Google Cloud Platform Compute Image and then use the GCP APIs to delete that image.
1.  [terraform_packer_example_test.go](/test/terraform_packer_example_test.go): Use Terratest to run `packer build` to
    build an AMI with a web server installed, deploy that AMI in AWS by running `terraform apply`, make HTTP requests to
    the web server to check that it is working correctly, and run `terraform destroy` to undeploy the web server.
1.  [terraform_gcp_example_test.go](/test/terraform_gcp_example_test.go): Use Terratest to run `terraform apply` on
    the Terraform GCP Example and verify you get the expected outputs.
1.  [terraform_remote_exec_example_test.go](/test/terraform_remote_exec_example_test.go): Use Terratest to run
    `terraform apply` and then remotely provision the instance while using a custom SSH agent managed by Terratest
1.  [terraform_scp_example_test.go](/test/terraform_scp_example_test.go): Use Terratest to simplify copying resources
    like config files and logs from deployed EC2 Instances. This is especially useful for getting a snapshot of the
    state of a deployment when a test fails.
1.  [kubernetes_basic_example_test.go](/test/kubernetes_basic_example_test.go): Use Terratest to run `kubectl apply`
    to apply a Kubernetes resource file, verify resources are created using the Kubernetes API, and then run `kubectl
    delete` to delete the resources at the end of the test.
1.  [kubernetes_rbac_example_test.go](/test/kubernetes_rbac_example_test.go): Use Terratest to run `kubectl apply` to
    apply a Kubernetes resource file, retrieve auth tokens to authenticate as the created ServiceAccount, update the
    kubeconfig file with the authentication token and add a new context to auth as the ServiceAccount, verify auth as
    the ServiceAccount by checking what resources you have access to, and finally run `kubectl delete` to delete the
    resources at the end of the test.
1.  [helm_basic_example_template_test.go](/test/helm_basic_example_template_test.go): Use Terratest to run `helm
    template` to test template rendering logic.


Finally, to see some real-world examples of Terratest in action, check out some of our open source infrastructure
modules:

1.  [Consul](https://github.com/hashicorp/terraform-aws-consul)
1.  [Vault](https://github.com/hashicorp/terraform-aws-vault)
1.  [Nomad](https://github.com/hashicorp/terraform-aws-nomad)
1.  [Couchbase](https://github.com/gruntwork-io/terraform-aws-couchbase/)
