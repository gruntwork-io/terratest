# External Consumer Simulation

This directory contains a test project that simulates how external consumers will use Terratest after modularization.

## Purpose

The consumer simulation validates that:
1. External users can import Terratest modules without ambiguous import errors
2. Dependencies resolve correctly across module boundaries
3. The modularization doesn't break the external API

## How It Works

This is a separate Go module that imports a representative mix of Terratest modules across different dependency tiers:
- **Tier 0**: `modules/testing` (no terratest dependencies)
- **Tier 1**: `modules/logger` (depends on testing)
- **Tier 3**: `modules/terragrunt` (depends on multiple modules)

## Testing

During **Phase 2** (workspace pilot), this project uses `replace` directives to point to local modules, similar to how the root module works.

During **Phase 3** (after releases), the `replace` directives will be removed, and this project will test that external consumers can use published module versions without workspace overrides.

## Running Tests

```bash
# With workspace (Phase 2)
cd test-external
go test -v ./...

# Without workspace (simulates external consumer)
cd test-external
GOWORK=off go test -v ./...
```

Both should pass during Phase 2, proving the modularization is working correctly.
