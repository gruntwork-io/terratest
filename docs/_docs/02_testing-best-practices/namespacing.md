---
layout: collection-browser-doc
title: Namespacing
category: testing-best-practices
excerpt: >-
  Learn how to avoid conflicts due to duplicated identifiers.
tags: ["testing-best-practices", "namespace", "id", "identifiers"]
order: 203
nav_title: Documentation
nav_title_link: /docs/
---

Just about all resources your tests create (e.g., servers, load balancers, machine images) should be "namespaced" with
a unique name to ensure that:

1.  You don't accidentally overwrite any "production" resources in that environment (though as mentioned in the previous
    section, your test environment should be completely isolated from prod anyway).
1.  You don't accidentally clash with other tests running in parallel.

For example, when deploying AWS infrastructure with Terraform, that typically means exposing variables that allow you
to configure auto scaling group names, security group names, IAM role names, and any other names that must be unique.

You can use Terratest's `random.UniqueId()` function to generate identifiers that are short enough to use in resource
names (just 6 characters) but random enough to make it unlikely that you'll have a conflict.

```go
uniqueId := random.UniqueId()
instanceName := fmt.Sprintf("terratest-http-example-%s", uniqueId)

terraformOptions := &terraform.Options {
  TerraformDir: "../examples/terraform-http-example",
  Vars: map[string]interface{} {
    "instance_name": instanceName,
  },
}

terraform.Apply(t, terraformOptions)
```
