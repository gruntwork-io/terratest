# Dynamically discover custom build tags from all .go files, excluding OS/arch/cgo tags
LINT_TAGS := $(shell grep -roh '//go:build .*' --include='*.go' . | sed 's|//go:build ||' | tr '|&()!' '\n' | tr -d ' ' | grep -E '^[a-z_]+$$' | sort -u | grep -Evx 'aix|android|darwin|dragonfly|freebsd|hurd|illumos|ios|js|linux|nacl|netbsd|openbsd|plan9|solaris|windows|zos|386|amd64|amd64p32|arm|armbe|arm64|arm64be|loong64|mips|mipsle|mips64|mips64le|mips64p32|mips64p32le|ppc|ppc64|ppc64le|riscv|riscv64|s390|s390x|sparc|sparc64|wasm|ignore|cgo' | tr '\n' ',' | sed 's/,$$//')

update-lint-config: SHELL:=/bin/bash
update-lint-config:
	curl -s https://raw.githubusercontent.com/gruntwork-io/terragrunt/main/.golangci.yml --output .golangci.yml
	tmpfile=$$(mktemp) ;\
	{ echo '# This file is generated from https://github.com/gruntwork-io/terragrunt/blob/main/.golangci.yml' ;\
	  echo '# It is automatically updated weekly via the update-lint-config workflow. Do not edit manually.' ;\
	  cat .golangci.yml; } > $${tmpfile} && mv $${tmpfile} .golangci.yml

lint:
	GOFLAGS="-tags=$(LINT_TAGS)" mise x golangci-lint -- golangci-lint run -v --timeout=30m ./...

lint-incremental:
	@echo "Incremental lint (new issues only)"
	GOFLAGS="-tags=$(LINT_TAGS)" mise x golangci-lint -- golangci-lint run -v --timeout=30m --new-from-merge-base=main ./...

lint-fix:
	@echo "Linting with auto-fix"
	GOFLAGS="-tags=$(LINT_TAGS)" mise x golangci-lint -- golangci-lint run -v --timeout=30m --fix ./...

.PHONY: lint lint-incremental lint-fix update-lint-config
