# Terragrunt Multi-Module Example

This folder contains a Terragrunt configuration with multiple modules that have dependencies (VPC → Database → App),
demonstrating how to use Terratest's `terragrunt` package to test multi-module setups.

Check out [modules/terragrunt/terragrunt_example_test.go](/modules/terragrunt/terragrunt_example_test.go) to see how you can write automated tests
for this configuration using `ApplyAll` and `DestroyAll`.

## Structure

```
.
├── modules/           # Terraform modules
│   ├── vpc/
│   ├── database/     # Depends on VPC
│   └── app/          # Depends on VPC and Database
└── live/             # Terragrunt configurations
    ├── vpc/
    ├── database/
    └── app/
```

## Running this module manually

1. Install [Terraform](https://www.terraform.io/) and [Terragrunt](https://terragrunt.gruntwork.io/).
1. `cd live`
1. Run `terragrunt apply --all`.
1. When you're done, run `terragrunt destroy --all`.

## Running automated tests against this module

1. Install [Terraform](https://www.terraform.io/) and [Terragrunt](https://terragrunt.gruntwork.io/).
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `go test -v -run TestTerragruntMultiModuleExample`
