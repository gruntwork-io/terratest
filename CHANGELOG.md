# Changelog

Per-submodule semver from v2 onward. See the [v2 migration guide](docs/_docs/04_migrating-to-v2/overview.md).

## Unreleased — v2.0.0

Splits Terratest into per-domain Go submodules. Runtime API unchanged.

- Every public module gets its own `go.mod` with the `/v2` SIV suffix (e.g. `modules/aws/v2`).
- Eight tier-0 utilities (collections, random, files, testing, environment, git, retry, logger) collapse into one `modules/core/v2` submodule as sub-packages.
- Minimum Go raised to 1.26 (1.21+ builds via Go's automatic toolchain).
- CI iterates per submodule for tidy, lint, and build. New enforcement: acyclic dep graph, single source per package path, `/v2` placement rule, no-replaces guard on release tags.
- `scripts/v2-import-rewrite.sh` for migrating consumer projects.
- `scripts/release-prep.sh` + `.github/workflows/v2-release.yml` automate the lockstep tag procedure.

### Backward compatibility

v1 (`github.com/gruntwork-io/terratest`) is frozen at its last tag, still served by `proxy.golang.org`. Pinned consumers are unaffected. v1 + v2 can coexist in the same `go.mod`. `go get -u` does not jump majors. Security backports land on `v1-maintenance`.

---

Prior history: see [GitHub Releases](https://github.com/gruntwork-io/terratest/releases) and `docs/_docs/03_migrating-to-v1/`.
