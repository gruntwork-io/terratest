# Terraform Database Example

This is a dummy database example which only contains necessary information of a database.

Check out [test/terraform_database_example_test.go](/test/terraform_database_example_test.go) to see how you can write automated tests for database. In order to make go test code work, you need to provide host, port, username, password and database name of a existing database, which you have already created on cloud platform or using docker before testing. Only Microsoft SQL Server, PostgreSQL and MySQL are supported.

## Running this module manually

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformDatabaseExample`