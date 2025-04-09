# Terraform Database Example

This example demonstrates how to create a PostgreSQL instance using Docker with Terraform.

Check out [test/terraform_database_example_test.go](/test/terraform_database_example_test.go) to learn how to write
automated tests for a database. To run the Go test code, you need to provide the host, port, username, password, and
database name of an existing database, which you should have already created on a cloud platform or using Docker before
running the tests. Currently, only Microsoft SQL Server, PostgreSQL, and MySQL are supported.

## Running this module manually

1. Install Terraform and ensure it's available in your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When finished, run `terraform destroy` to remove the resources.

## Running automated tests against this module

1. Install [Terraform](https://www.terraform.io/) and ensure it's available in your `PATH`.
1. Install [Golang](https://golang.org/) and ensure this code is checked out into your `GOPATH`.
1. Run `go mod tidy` to manage dependencies.
1. Run `go test -v test/terraform_database_example_test.go` to execute the tests.