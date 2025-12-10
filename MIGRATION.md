# Migration Guide: Terratest v1.0.0

Starting with v1.0.0, Terratest is split into independent Go modules.

## Migration

**Before:**
```bash
go get github.com/gruntwork-io/terratest
```

**After:**
```bash
go get github.com/gruntwork-io/terratest/modules/terraform@v1.0.0
go get github.com/gruntwork-io/terratest/modules/aws@v1.0.0
```

Import paths remain unchanged. Run `go mod tidy` after updating.

## Available Modules

All modules are under `modules/`: terraform, aws, azure, gcp, k8s, helm, docker, terragrunt, ssh, shell, logger, retry, random, files, http-helper, dns-helper, test-structure, collections, environment, git, oci, opa, packer, database, slack, version-checker, testing.
