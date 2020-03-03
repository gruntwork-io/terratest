---
layout: collection-browser-doc
title: Testing plan output
category: testing-best-practices
excerpt: >-
  Testing specific changes described in Terraform plan output
tags: ["testing-best-practices"]
order: 202
nav_title: Documentation
nav_title_link: /docs/
---

TODO: change the `order`s in the preambles of all of these best practices pages

Sometimes you can get the information you need in your tests from a `terraform plan`; you may not need to do an `apply`.  Terratest provides a way to use certain information from a `terraform plan` in your tests.

The `InitAndPlanWithInfo` function will return a `PlanInfo` struct that will contain information about which resources would be changed and what specifically about them would change.

For example, let's say you have some `web_server` resources that you're creating with [`for_each`](https://www.terraform.io/docs/configuration/resources.html#for_each-multiple-resource-instances-defined-by-a-map-or-set-of-strings), and you're using the `for_each` values in the name:

```hcl
resource "aws_instance" "web_server" {
  for_each = toset(["env1", "env2", "env3"])
  name = "web_server_${each.value}"
  ...
}
```

if you want to assert that the names of `web_server` resources that are defined with a terraform `for_each`, you can test this by making assertions against info in this `PlanInfo` struct:

```go
plan := terraform.InitAndPlanWithInfo(t, terraformOptions)

found := false
for _, pv := range plan.PlannedValues {
    if pv.Name == "web_server" {
        expected := fmt.Sprintf("web_server_%s", rc.Index)
        assert.Equal(rc.Values["name"], expected)
        found = true
    }
}
assert.True(found)
```

For a full description of what Terraform plan information is available in the PlanInfo struct, see TODO.
