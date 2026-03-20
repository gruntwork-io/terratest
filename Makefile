update-lint-config: SHELL:=/bin/bash
update-lint-config:
	curl -s https://raw.githubusercontent.com/gruntwork-io/terragrunt/main/.golangci.yml --output .golangci.yml
	tmpfile=$$(mktemp) ;\
	{ echo '# This file is generated from https://github.com/gruntwork-io/terragrunt/blob/main/.golangci.yml' ;\
	  echo '# It is automatically updated weekly via the update-lint-config workflow. Do not edit manually.' ;\
	  cat .golangci.yml; } > $${tmpfile} && mv $${tmpfile} .golangci.yml

lint:
	golangci-lint run ./...

lint-allow-list:
	golangci-lint run ./modules/random/... ./modules/testing/... ./modules/slack/... ./modules/collections/... ./modules/environment/... ./modules/retry/... ./modules/shell/... ./modules/git/... ./modules/files/... ./modules/oci/...

.PHONY: lint update-lint-config
