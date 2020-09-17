---
layout: collection-browser-doc
title: Timeouts and logging
category: testing-best-practices
excerpt: >-
  Long-running infrastructure tests may exceed timeouts or can be killed if they do not prompt logs.
tags: ["testing-best-practices", "timeout", "error"]
order: 205
nav_title: Documentation
nav_title_link: /docs/
---

Go's package testing has a default timeout of 10 minutes, after which it forcibly kills your tests—even your cleanup
code won't run! It's not uncommon for infrastructure tests to take longer than 10 minutes, so you'll almost always
want to increase the timeout by using the `-timeout` option, which takes a `go` duration string (e.g `10m` for 10
minutes or `1h` for 1 hour):

```bash
go test -timeout 30m
```

Note that many CI systems will also kill your tests if they don't see any log output for a certain period of time
(e.g., 10 minutes in CircleCI). If you use Go's `t.Log` and `t.Logf` for logging in your tests, you'll find that these
functions buffer all log output until the very end of the test (see https://github.com/golang/go/issues/24929 for more
info). If you have a long-running test, this might mean you get no log output for more than 10 minutes, and the CI
system will shut down your tests. Moreover, if your test has a bug that causes it to hang, you won't see any log output
at all to help you debug it.

Therefore, we recommend instead using Terratest's `logger.Log` and `logger.Logf` functions, which log to `stdout`
immediately:

```go
func TestFoo(t *testing.T) {
  logger.Log(t, "This will show up in stdout immediately")
}
```

Finally, if you're testing multiple Go packages, be aware that Go will buffer log output—even that sent directly to
`stdout` by `logger.Log` and `logger.Logf`—until all the tests in the package are done. This leads to the same
difficulties with CI servers and debugging. The workaround is to tell Go to test each package sequentially using the
`-p 1` flag:

```bash
go test -timeout 30m -p 1 ./...
```

See the [Cleanup]({{site.baseurl}}/docs/testing-best-practices/cleanup/) for more information on how to setup robust clean up procedures in the face of test timeouts and instabilities.
