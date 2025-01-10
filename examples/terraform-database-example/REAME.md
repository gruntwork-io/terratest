# Terraform Database Example

This creates an example postgres instance using docker.

Check out [test/terraform_database_example_test.go](/test/terraform_database_example_test.go) to see how you can write automated tests for database. In order to make go test code work, you need to provide host, port, username, password and database name of a existing database, which you have already created on cloud platform or using docker before testing. Only Microsoft SQL Server, PostgreSQL and MySQL are supported.

## Running this module manually

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
2. Run `terraform init`.
3. Run `terraform apply`.
4. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
2. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
3. `go mod tidy`
4. `go test -v test/terraform_database_example_test.go`
