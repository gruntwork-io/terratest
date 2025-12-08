# Migration Guide: Terratest Modularization (v0.55.0+)

## What Changed

Starting with v0.55.0, Terratest has been split into independent Go modules. Import only the modules you need instead of the entire library.

## Migration Steps

### 1. Update your go.mod

**Before:**
```bash
go get github.com/gruntwork-io/terratest
```

**After:**
```bash
go get github.com/gruntwork-io/terratest/modules/terraform@v0.55.0
go get github.com/gruntwork-io/terratest/modules/aws@v0.55.0
# Add other modules as needed
```

### 2. Import paths stay the same

Your code doesn't change:
```go
import (
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/gruntwork-io/terratest/modules/aws"
)
```

### 3. Clean up

```bash
go mod tidy
```

## Available Modules

| Module | Path |
|--------|------|
| terraform | `modules/terraform` |
| aws | `modules/aws` |
| azure | `modules/azure` |
| gcp | `modules/gcp` |
| k8s | `modules/k8s` |
| helm | `modules/helm` |
| docker | `modules/docker` |
| terragrunt | `modules/terragrunt` |
| ssh | `modules/ssh` |
| shell | `modules/shell` |
| logger | `modules/logger` |
| retry | `modules/retry` |
| random | `modules/random` |
| files | `modules/files` |
| http-helper | `modules/http-helper` |
| dns-helper | `modules/dns-helper` |
| test-structure | `modules/test-structure` |
| collections | `modules/collections` |
| environment | `modules/environment` |
| git | `modules/git` |
| oci | `modules/oci` |
| opa | `modules/opa` |
| packer | `modules/packer` |
| database | `modules/database` |
| slack | `modules/slack` |
| version-checker | `modules/version-checker` |
| testing | `modules/testing` |

## Troubleshooting

**"module does not contain package" error:**
The root module no longer exists. Use a specific submodule path (e.g., `modules/terraform`).

**"ambiguous import" error:**
Run `go mod tidy` to clean up dependencies.
