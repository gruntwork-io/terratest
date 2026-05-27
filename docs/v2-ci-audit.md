# Terratest v2 CI Audit

Findings from auditing existing GitHub Actions workflows against the planned per-submodule `go.mod` layout. PR #1632 reportedly "never got past CI setup" — these are the most likely reasons why.

## Breakage points

### 1. `.github/workflows/go-mod-tidy-check.yml` — BREAKS

Runs `go mod tidy` and `git diff go.mod go.sum` at repo root only. Will silently miss untidy state in `modules/<name>/go.mod` files.

**Fix:** iterate.

```yaml
- name: Run go mod tidy in all submodules
  run: |
    set -e
    find . -name go.mod -not -path '*/node_modules/*' -print0 \
      | xargs -0 -n1 dirname \
      | while read dir; do
          echo "Tidying $dir"
          (cd "$dir" && go mod tidy)
        done

- name: Check for changes
  run: |
    if ! git diff --exit-code '**/go.mod' '**/go.sum'; then
      echo "::error::Run 'go mod tidy' in the modules with diffs."
      exit 1
    fi
```

### 2. `.github/workflows/lint.yml` + `Makefile` `lint` target — BREAKS

`make lint` runs `golangci-lint run -v ./...` from repo root. `./...` is module-scoped — it only lints packages in the root module. Submodules get skipped entirely.

**Fix options:**

a. **Use `go.work`** at repo root listing all submodules. Recent golangci-lint (v1.55+) respects go.work and lints across the workspace. Cleanest if the version is recent enough. Verify with `golangci-lint version` in CI.

b. **Iterate via Makefile.** Loop over `find . -name go.mod` and run lint per submodule. Slower (cold cache per run) but unambiguous.

Recommendation: (a) if golangci-lint version allows, (b) as fallback.

### 3. `.github/workflows/build-and-release.yml` — BREAKS

```yaml
go build -o "cmd/bin/terratest_log_parser_..." ./cmd/terratest_log_parser
```

After OSS-3534 splits `cmd/terratest_log_parser` into its own go.mod, this path is no longer a package in the root module. Go will error with "no required module provides package".

**Fix:** `cd` into the cmd directory.

```yaml
- name: Build binaries
  run: |
    for cmd_dir in cmd/*/; do
      cmd_name=$(basename "$cmd_dir")
      (cd "$cmd_dir" && CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH \
        go build -o "../../cmd/bin/${cmd_name}_${OS}_${ARCH}${EXT}" .)
    done
```

### 4. `.github/workflows/{aws,azure,gcp,k8s,terraform,terragrunt}-*-tests.yml` — DESIGN QUESTION

All integration test workflows run `go test … ./test/...` from the root. The `test/` directory currently lives in the root module and imports `modules/<various>`. After modularization, `test/` files will import multiple submodules.

Three options:

a. **`test/` stays in the root module.** Root `go.mod` requires tagged versions of all submodules. Local changes in `modules/<name>` are NOT seen by `go test ./test/...` unless you tag-and-replace, or add `go.work`.

b. **`go.work` at repo root** including `.`, `./modules/*`, `./cmd/*`. Then `go test ./test/...` resolves submodules via the workspace. Works for CI and local dev. Not published to consumers (go.work isn't tagged).

c. **`test/` becomes its own go.mod** with local replaces for everything. Heavy.

Recommendation: (b). `go.work` is the standard answer for multi-module monorepos. Add `go.work` to repo root, add `go.work.sum` to .gitignore (or commit it — consensus debated).

## Other observations

### Pre-commit (`pre-commit-config.yaml` + workflow)

The hooks (`goimports`, `terraform-fmt`, `test-interfaces-used`) are file-based and don't read `go.mod`. Should work unchanged. The `test-interfaces-used` hook greps for `*testing.T` in `modules/**/*.go` — still works regardless of module boundaries.

### Mise + tool versions

Both `mise.toml` (root) and `mise-action` in workflows. Tools (`go`, `golangci-lint`, etc.) are pinned at the repo level. Doesn't break with multi-module — same toolchain applies to all submodules. Each submodule's `go.mod` `go 1.26` line must match the mise-installed Go version. (Currently all submodules use `go 1.26`.)

### Caching

Multiple workflows cache `~/go/pkg/mod` keyed on `hashFiles('**/go.sum')`. The `**/go.sum` glob will pick up submodule go.sum files too, so the cache key naturally invalidates when any submodule changes. Cache *fill* will be larger (more modules → more downloaded artifacts), but no correctness issue.

## Minimal CI pre-flight before OSS-3526 starts

1. Add `go.work` to repo root listing the modules already split in bundle 1.
2. Update `go-mod-tidy-check.yml` to iterate.
3. Update `Makefile` `lint` target to either use the go.work workspace or iterate.
4. Update `build-and-release.yml` to `cd` per `cmd/`.
5. Verify a no-op PR passes CI before any real submodule PR lands.

Doing these in a single "v2-ci-prep" PR before OSS-3526 is much safer than discovering them mid-stack like PR #1632 did.

## Risk that doesn't show up here

`go.work` is not published; downstream consumers of terratest v2 will not see it. So CI green with go.work + local changes does NOT prove consumers can fetch the modules. The lockstep release runbook's post-tag verification (a scratch consumer hitting proxy.golang.org) is the only check that catches consumer-side breakage. Both checks are needed; one doesn't substitute for the other.
