# Terratest v2 Lockstep Release Runbook

This runbook covers tagging all v2 submodules (≈12) in a coordinated release.

**Critical constraint:** `proxy.golang.org` caches every published `(module, version)` pair immutably. A bad tag cannot be fixed by retagging — it must be superseded by a new patch version. There is no rollback.

## Pre-flight checks

Run these on the release branch BEFORE pushing any tag:

1. **All submodules build cleanly from their own root.**
   For each `modules/<name>/go.mod`:
   ```
   (cd modules/<name> && go build ./... && go vet ./...)
   ```

2. **All tests pass.** Per submodule:
   ```
   (cd modules/<name> && go test -count=1 ./...)
   ```

3. **No local `replace` directives in any submodule `go.mod`** (other than dev-only `replace` for the main repo's own modules — those become pinned tagged versions at release).
   ```
   grep -nH '^replace' modules/*/go.mod
   ```
   The release version of each `go.mod` must `require` tagged versions of sibling submodules, NOT point at `../<name>`. Strip the local replaces in the release commit.

4. **`go.sum` files committed and consistent.**

5. **Release branch fully merged from `main`.** Tag from the merge commit, not a feature branch.

6. **Dry-run with a `-rc.N` tag set** on a fork to verify the same matrix in a controlled namespace before touching the canonical repo.

## Tag ordering (dependency-first)

Tags must be pushed in dependency order. If module B requires module A at version `vX`, A's tag MUST exist on the proxy before B's go.mod is published — otherwise consumers of B will hit "unknown revision" errors on first fetch.

Order based on the dependency tree:

1. `modules/core/v2.0.0` — utility packages (collections, random, files, testing, environment, git, retry, logger collapsed in)
2. `modules/shell/v2.0.0`, `modules/http-helper/v2.0.0`, `modules/dns-helper/v2.0.0`, `modules/ssh/v2.0.0`, `modules/version-checker/v2.0.0` — mid-stack helpers (depend on core)
3. `modules/docker/v2.0.0`, `modules/packer/v2.0.0`, `modules/database/v2.0.0`, `modules/slack/v2.0.0`, `modules/oci/v2.0.0`, `modules/opa/v2.0.0` — mid-stack tooling (depend on core + helpers)
4. `modules/k8s/v2.0.0`, `modules/helm/v2.0.0` — Kubernetes (depend on core + tooling)
5. `modules/aws/v2.0.0`, `modules/azure/v2.0.0`, `modules/gcp/v2.0.0` — clouds (depend on core + tooling)
6. `modules/terraform/v2.0.0`, `modules/terragrunt/v2.0.0`, `modules/test-structure/v2.0.0` — Terraform layer
7. `cmd/<each>/v2.0.0` — binaries last (they consume everything)

Each tier must be fully pushed and proxy-verified (see post-tag checks) before starting the next.

## Tag command shape

```
git tag -a modules/<name>/v2.0.0 -m "v2.0.0 release"
git push origin modules/<name>/v2.0.0
```

All tags point at the **same commit SHA** — the release commit on the release branch.

## Post-tag verification (per submodule)

After each push, verify the proxy serves it before moving to the next tier:

```
for m in <names>; do
  echo -n "$m: "
  curl -s -o /dev/null -w '%{http_code}\n' \
    "https://proxy.golang.org/github.com/gruntwork-io/terratest/modules/$m/v2/@v/v2.0.0.info"
done
```

All 200s. A non-200 means the proxy hasn't caught up yet (rare, but wait 60s and retry) OR the tag format is wrong (stop, investigate).

End-to-end:
```
mkdir /tmp/v2-smoke && cd /tmp/v2-smoke
go mod init smoke
go get github.com/gruntwork-io/terratest/modules/core/v2@v2.0.0
go get github.com/gruntwork-io/terratest/modules/aws/v2@v2.0.0
# ... write a tiny main.go importing each, build it
go build ./...
```

## Partial-failure recovery

**If a tag is pushed with a bug:**
- The tag is immutable on the proxy. You cannot fix it in place.
- Cut a new patch: `v2.0.1`. All later modules in the chain that already depended on `v2.0.0` will need their `require` lines bumped to `v2.0.1` and themselves re-tagged.
- This cascades. Pre-flight checks exist to prevent this — do not skip them.

**If a tag push fails mid-chain:**
- Inspect the failed push (network, auth, GitHub rate limit). Retry the same tag command — git tag push is idempotent.
- DO NOT delete and recreate the tag locally if the push partially succeeded. Verify with `git ls-remote --tags origin | grep <name>`.

**If the proxy serves a wrong commit:**
- Should not happen — tags are by SHA. If it does, file a bug with the Go team. Workaround: cut a new patch tag on a new commit.

## Dry-run procedure

Before the real release:

1. Branch off `main` to a release-prep branch.
2. Strip local `replace` directives, pin tagged versions in each `go.mod`.
3. Push to a fork (e.g., `<user>/terratest`).
4. Run the full tag matrix on the fork using `v2.0.0-rc.1` versions.
5. Run the post-tag verification (curl + scratch consumer).
6. If all green, repeat on the canonical repo with `v2.0.0` tags.

## Roles and sequencing

- **Release engineer**: pushes tags, runs verification.
- **Second pair of eyes**: reviews the strip-`replace` commit before tag.
- **Comms**: announcement post (blog + CHANGELOG) sent only AFTER all tags are verified green.

## Reference

- Module versioning rules: https://go.dev/ref/mod#major-version-suffixes
- Proxy caching behavior: https://proxy.golang.org/
- SIV requirement for v2+: https://go.dev/blog/v2-go-modules
