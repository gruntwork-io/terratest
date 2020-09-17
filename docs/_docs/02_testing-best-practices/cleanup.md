---
layout: collection-browser-doc
title: Cleanup
category: testing-best-practices
excerpt: >-
  Since automated tests with Terratest deploy real resources into real environments, you'll want to make sure your tests
  always cleanup after themselves.
tags: ["testing-best-practices", "clean", "terraform-destroy", "terraform-apply"]
order: 204
nav_title: Documentation
nav_title_link: /docs/
---

Since automated tests with Terratest deploy real resources into real environments, you'll want to make sure your tests
always cleanup after themselves so you don't leave a bunch of resources lying around. Typically, you should use Go's
`defer` keyword to ensure that the cleanup code always runs, even if the test hits an error along the way.

For example, if your test runs `terraform apply`, you should run `terraform destroy` at the end to clean up:

```go
// Ensure cleanup always runs
defer terraform.Destroy(t, options)

// Deploy
terraform.Apply(t, options)

// Validate
checkServerWorks(t, options)
```

Of course, despite your best efforts, occasionally cleanup will fail, perhaps due to the CI server going down, or a bug
in your code, or a temporary network outage. To handle those cases, we run a tool called
[cloud-nuke](https://github.com/gruntwork-io/cloud-nuke) in our test AWS account on a nightly basis to clean up any
leftover resources.
