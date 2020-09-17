---
layout: collection-browser-doc
title: Avoid test caching
category: testing-best-practices
excerpt: >-
  Since Go 1.10, test results are automatically cached. See how to turn off caching test results.
tags: ["testing-best-practices", "cache"]
order: 207
nav_title: Documentation
nav_title_link: /docs/
---

Since Go 1.10, test results are automatically [cached](https://golang.org/doc/go1.10#test). This can lead to Go not
running your tests again if you haven't changed any of the Go code. Since you're probably mainly manipulating Terraform
files, you should consider turning the caching of test results off. This ensures that the tests are run every time
you run `go test` and the result is not just read from the cache.

To turn caching off, you can set the `-count` flag to `1` force the tests to run:

```shell
$ go test -count=1 -timeout 30m -p 1 ./...
```
