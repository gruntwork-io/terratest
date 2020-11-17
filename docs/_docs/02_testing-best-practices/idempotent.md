---
layout: collection-browser-doc
title: Idempotent
category: testing-best-practices
excerpt: >-
  Test that your Terraform configuration results in consistent deployments.
tags: ["testing-best-practices", "idempotent", "terraform"]
order: 212
nav_title: Documentation
nav_title_link: /docs/
---

A Terraform configuration is idempotent when a second apply results in 0 changes. An idempotent configuration ensures that:

1.  What you define in Terraform is exactly what is being deployed. 
1.  Detection of bugs in Terraform resources and providers that might affect your configuration.

You can use Terratest's `terraform.ApplyAndIdempotent()` function to both apply your Terraform configuration and test its
idempotency.

```go
terraform.ApplyAndIdempotent(t, terraformOptions)
```

If a second apply of your Terraform configuration results in changes then your test will fail.
