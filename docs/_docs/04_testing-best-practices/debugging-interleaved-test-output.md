---
layout: collection-browser-doc
title: Debugging interleaved test output
category: testing-best-practices
excerpt: >-
  Learn more about ...
tags: ["testing-best-practices"]
order: 406
nav_title: Documentation
nav_title_link: /docs/
---

**Note**: The `terratest_log_parser` requires an explicit installation. See [Installing the utility
binaries](#installing-the-utility-binaries) for installation instructions.

If you log using Terratest's `logger` package, you may notice that all the test outputs are interleaved from the
parallel execution. This may make it difficult to debug failures, as it can be tedious to sift through the logs to find
the relevant entries for a failing test, let alone find the test that failed.

Therefore, Terratest ships with a utility binary `terratest_log_parser` that can be used to break out the logs.

To use the utility, you simply give it the log output from a `go test` run and a desired output directory:

```bash
go test -timeout 30m | tee test_output.log
terratest_log_parser -testlog test_output.log -outputdir test_output
```

This will:

- Create a file `TEST_NAME.log` for each test it finds from the test output containing the logs corresponding to that
  test.
- Create a `summary.log` file containing the test result lines for each test.
- Create a `report.xml` file containing a Junit XML file of the test summary (so it can be integrated in your CI).

The output can be integrated in your CI engine to further enhance the debugging experience. See Terratest's own
[circleci configuration](/.circleci/config.yml) for an example of how to integrate the utility with CircleCI. This
provides for each build:

- A test summary view showing you which tests failed:

![CircleCI test summary](/_docs/images/circleci-test-summary.png)

- A snapshot of all the logs broken out by test:

![CircleCI logs](/_docs/images/circleci-logs.png)
