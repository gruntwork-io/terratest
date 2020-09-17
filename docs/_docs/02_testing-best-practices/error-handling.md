---
layout: collection-browser-doc
title: Error handling
category: testing-best-practices
excerpt: >-
  Learn how to handle errors.
tags: ["testing-best-practices", "terraform", "error"]
order: 208
nav_title: Documentation
nav_title_link: /docs/
---

Just about every method `foo` in Terratest comes in two versions: `foo` and `fooE` (e.g., `terraform.Apply` and
`terraform.ApplyE`).

- `foo`: The base method takes a `t *testing.T` as an argument. If the method hits any errors, it calls `t.Fatal` to
  fail the test.

- `fooE`: Methods that end with the capital letter `E` always return an `error` as the last argument and never call
  `t.Fatal` themselves. This allows you to decide how to handle errors.

You will use the base method name most of the time, as it allows you to keep your code more concise by avoiding
`if err != nil` checks all over the place:

```go
terraform.Init(t, terraformOptions)
terraform.Apply(t, terraformOptions)
url := terraform.Output(t, terraformOptions, "url")
```

In the code above, if `Init`, `Apply`, or `Output` hits an error, the method will call `t.Fatal` and fail the test
immediately, which is typically the behavior you want. However, if you are _expecting_ an error and don't want it to
cause a test failure, use the method name that ends with a capital `E`:

```go
if _, err := terraform.InitE(t, terraformOptions); err != nil {
  // Do something with err
}

if _, err := terraform.ApplyE(t, terraformOptions); err != nil {
  // Do something with err
}

url, err := terraform.OutputE(t, terraformOptions, "url")
if err != nil {
  // Do something with err
}
```

As you can see, the code above is more verbose, but gives you more flexibility with how to handle errors.
