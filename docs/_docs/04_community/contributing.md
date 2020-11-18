---
layout: collection-browser-doc
title: Contributing
category: community
excerpt: >-
  Terratest is an open source project, and contributions from the community are very welcome!
tags: ["contributing", "community"]
order: 400
nav_title: Documentation
nav_title_link: /docs/
---

Terratest is an open source project, and contributions from the community are very welcome\! Please check out the
[Contribution Guidelines](#contribution-guidelines) and [Developing Terratest](#developing-terratest) for
instructions.

## Contribution Guidelines

Contributions to this repo are very welcome! We follow a fairly standard [pull request
process](https://help.github.com/articles/about-pull-requests/) for contributions, subject to the following guidelines:

1. [Types of contributions](#types-of-contributions)
1. [File a GitHub issue](#file-a-github-issue)
1. [Update the documentation](#update-the-documentation)
1. [Update the tests](#update-the-tests)
1. [Update the code](#update-the-code)
1. [Create a pull request](#create-a-pull-request)
1. [Merge and release](#merge-and-release)

### Types of contributions

Broadly speaking, Terratest contains two types of helper functions:

1. Integrations with external tools
1. Infrastructure and validation helpers

We accept different types of contributions for each of these two types of helper functions, as described next.

#### Integrations with external tools

These are helper functions that integrate with various DevOps tools—e.g., Terraform, Docker, Packer, and
Kubernetes—that you can use to deploy infrastructure in your automated tests. Examples:

* `terraform.InitAndApply`: run `terraform init` and `terraform apply`.
* `packer.BuildArtifacts`: run `packer build`.
* `shell.RunCommandAndGetOutput`: run an arbitrary shell command and return `stdout` and `stderr` as a string.

Here are the guidelines for contributions with external tools:

1. **Fixes and improvements to existing integrations**: All bug fixes and new features for existing tool integrations
   are very welcome!  

1. **New integrations**: Before contributing an integration with a totally new tool, please file a GitHub issue to
   discuss with us if it's something we are interested in supporting and maintaining. For example, we may be open to
   new integrations with Docker and Kubernetes tools, but we may not be open to integrations with Chef or Puppet, as
   there are already testing tools available for them.

#### Infrastructure and validation helpers

These are helper functions for creating, destroying, and validating infrastructure directly via API calls or SDKs.
Examples:

* `http_helper.HttpGetWithRetry`: make an HTTP request, retrying until you get a certain expected response.
* `ssh.CheckSshCommand`: SSH to a server and execute a command.
* `aws.CreateS3Bucket`: create an S3 bucket.
* `aws.GetPrivateIpsOfEc2Instances`:  use the AWS APIs to fetch IPs of some EC2 instances.

The number of possible such helpers is nearly infinite, so to avoid Terratest becoming a gigantic, sprawling library
we ask that contributions for new infrastructure helpers are limited to:

1. **Platforms**: we currently only support three major public clouds (AWS, GCP, Azure) and Kubernetes. There is some
   code contributed earlier for other platforms (e.g., OCI), but until we have the time/resources to support those
   platforms fully, we will only accept contributions for the major public clouds and Kubernetes.

1. **Complexity**: we ask that you only contribute infrastructure and validation helpers for code that is relatively
   complex to do from scratch. For example, a helper that merely wraps an existing function in the AWS or GCP SDK is
   not a great choice, as the wrapper isn't contributing much value, but is bloating the Terratest API. On the other
   hand, helpers that expose simple APIs for complex logic are great contributions: `ssh.CheckSshCommand` is a great
   example of this, as it provides a simple one-line interface for dozens of lines of complicated SSH logic.

1. **Popularity**: Terratest should only contain helpers for common use cases that come up again and again in the
   course of testing. We don't want to bloat the library with lots of esoteric helpers for rarely used tools, so
   here's a quick litmus test: (a) Is this helper something you've used once or twice in your own tests, or is it
   something you're using over and over again? (b) Does this helper only apply to some use case specific to your
   company or is it likely that many other Terratest users are hitting this use case over and over again too?

1. **Creating infrastructure**: we try to keep helper functions that create infrastructure (e.g., use the AWS SDK to
   create an S3 bucket or EC2 instance) to a minimum, as those functions typically require maintaining state (so that
   they are idempotent and can clean up that infrastructure at the end of the test) and dealing with asynchronous and
   eventually consistent cloud APIs. This can be surprisingly complicated, so we typically recommend using a tool like
   Terraform, which already handles all that complexity, to create any infrastructure you need at test time, and
   running Terratest's built-in `terraform` helpers as necessary. If you're considering contributing a function that
   creates infrastructure directly (e.g., using a cloud provider's APIs), please file a GitHub issue to explain why
   such a function would be a better choice than using a tool like Terraform.

### File a GitHub issue

Before starting any work, we recommend filing a GitHub issue in this repo. This is your chance to ask questions and
get feedback from the maintainers and the community before you sink a lot of time into writing (possibly the wrong)
code. If there is anything you're unsure about, just ask!

### Update the documentation

We recommend updating the documentation *before* updating any code (see [Readme Driven
Development](http://tom.preston-werner.com/2010/08/23/readme-driven-development.html)). This ensures the documentation
stays up to date and allows you to think through the problem at a high level before you get lost in the weeds of
coding.

The documentation is built with Jekyll and hosted on the Github Pages from `docs` folder on `master` branch. Check out [Terratest website](https://github.com/gruntwork-io/terratest/tree/master/docs#working-with-the-documentation) to learn more about working with the documentation.

### Update the tests

We also recommend updating the automated tests *before* updating any code (see [Test Driven
Development](https://en.wikipedia.org/wiki/Test-driven_development)). That means you add or update a test case,
verify that it's failing with a clear error message, and *then* make the code changes to get that test to pass. This
ensures the tests stay up to date and verify all the functionality in this Module, including whatever new
functionality you're adding in your contribution. The instructions for running the automated tests can be
found [here](https://terratest.gruntwork.io/docs/community/contributing/#developing-terratest).

### Update the code

At this point, make your code changes and use your new test case to verify that everything is working. As you work,
please make every effort to avoid unnecessary backwards incompatible changes. This generally means that you should
not delete or rename anything in a public API.

If a backwards incompatible change cannot be avoided, please make sure to call that out when you submit a pull request,
explaining why the change is absolutely necessary.

Note that we use pre-commit hooks with this project. To ensure they run:

1. Install [pre-commit](https://pre-commit.com/).
1. Run `pre-commit install`.

One of the pre-commit hooks we run is [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports). To prevent the
hook from failing, make sure to :

1. Install [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
1. Run `goimports -w .`.

We have a [style guide](https://gruntwork.io/guides/style%20guides/golang-style-guide/) for the Go programming language,
in which we documented some best practices for writing Go code. Please ensure your code adheres to the guidelines
outlined in the guide.

### Create a pull request

[Create a pull request](https://help.github.com/articles/creating-a-pull-request/) with your changes. Please make sure
to include the following:

1. A description of the change, including a link to your GitHub issue.
1. The output of your automated test run, preferably in a [GitHub Gist](https://gist.github.com/). We cannot run
   automated tests for pull requests automatically due to [security
   concerns](https://circleci.com/docs/2.0/oss/#security), so we need you to manually provide this
   test output so we can verify that everything is working.
1. Any notes on backwards incompatibility or downtime.

#### Validate the Pull Request for Azure Platform

If you're contributing code for the [Azure Platform](https://azure.com) and if you have an active _Azure subscription_, it's recommended to follow the below guidelines after [creating a pull request](https://help.github.com/articles/creating-a-pull-request/). If you're contributing code for any other platform (e.g., AWS, GCP, etc), you can skip these steps.

> Once the Terratest maintainers add `Azure` tag and _Approve_ the PR, following pipeline will run automatically to perform a full validation of the Azure contribution. You also can run the pipeline manually on your forked repo by following the below guideline.


We have a separate CI pipeline for _Azure_ code. To run it on a forked repo:

1. Run the following [Azure Cli](https://docs.microsoft.com/cli/azure/) command on your preferred Terminal to create Azure credentials and copy the output:

    ```bash
    az ad sp create-for-rbac --name "terratest-az-cli" --role contributor --sdk-auth
    ```

1. Go to Secrets settings page under `Settings` tab in your forked project, `https://github.com/<YOUR_GITHUB_ACCOUNT>/terratest/settings`, on GitHub.

1. Create a new `Secret` named `AZURE_CREDENTIALS` and paste the Azure credentials you copied from the 1<sup>st</sup> step as the value

    > `AZURE_CREDENTIALS` will be stored in _your_ GitHub account; neither the Terratest maintainers nor anyone else will have any access to it. Under the hood, GitHub stores your secrets in a secure, encrypted format (see: [GitHub Actions Secrets Reference](https://docs.github.com/en/free-pro-team@latest/actions/reference/encrypted-secrets) for more information). Once the secret is created, it's only possible to update or delete it; the value of the secret can't be viewed. GitHub uses a [libsodium sealed box](https://libsodium.gitbook.io/doc/public-key_cryptography/sealed_boxes) to help ensure that secrets are encrypted before they reach GitHub.

1. Create a [new Personal Access Token (PAT)](https://github.com/settings/tokens/new) page under [Settings](https://github.com/settings/profile) / [Developer Settings](https://github.com/settings/apps), making sure `write:discussion` and `public_repo` scopes are checked. Click the _Generate token_ button and copy the generated PAT.

1. Go back to settings/secrets in your fork and [Create a new Secret](https://docs.github.com/actions/reference/encrypted-secrets#creating-encrypted-secrets-for-a-repository) named `PAT`.  Paste the output from the 4<sup>th</sup> step as the value

    > `PAT` will be stored in _your_ GitHub account; neither the Terratest maintainers nor anyone else will have any access to it. Under the hood, GitHub stores your secrets in a secure, encrypted format (see: [GitHub Actions Secrets Reference](https://docs.github.com/en/free-pro-team@latest/actions/reference/encrypted-secrets) for more information). Once the secret is created, it's only possible to update or delete it; the value of the secret can't be viewed. GitHub uses a [libsodium sealed box](https://libsodium.gitbook.io/doc/public-key_cryptography/sealed_boxes) to help ensure that secrets are encrypted before they reach GitHub.

1. Go to Actions tab on GitHub (https://github.com/<GITHUB_ACCOUNT>/terratest/actions)

1. Click `ci-workflow` workflow

1. Click `Run workflow` button and fill the fields in the drop down
    * _Repository Info_ : name of the forked repo (_e.g. xyz/terratest_)
    * _Name of the branch_ : branch name on the forked repo (_e.g. feature/adding-some-important-module_)
    * _Name of the official terratest repo_ : home of the target pr (_gruntwork-io/terratest_)
    * PR number on the official terratest repo : pr number on the official terratest repo (_e.g. 14, 25, etc._).  Setting this value will leave a success/failure comment in the PR once CI completes execution.

    * Skip provider registration : set true if you want to skip terraform provider registration for debug purposes (_false_ or _true_)

1. Wait for the `ci-workflow` to be finished

    > The pipeline will use the given Azure subscription and deploy real resources in your Azure account as part of running the test. When the tests finish, they will tear down the resources they created. Of course, if there is a bug or glitch that prevents the clean up code from running, some resources may be left behind, but this is rare. Note that these resources may cost you money! You are responsible for all charges in your Azure subscription.

1. PR with the given _PR Number_ will have the result of the `ci-workflow` as a comment

### Merge and release

The maintainers for this repo will review your code and provide feedback. Once the PR is accepted, they will merge the
code and release a new version, which you'll be able to find in the [releases page](https://github.com/gruntwork-io/terratest/releases).

## Developing Terratest

1. [Running tests](#running-tests)
1. [Versioning](#versioning)

### Running tests

Terratest itself includes a number of automated tests.

**Note #1**: Some of these tests create real resources in an AWS account. That means they cost money to run, especially
if you don't clean up after yourself. Please be considerate of the resources you create and take extra care to clean
everything up when you're done!

**Note #2**: In order to run tests that access your AWS account, you will need to configure your [AWS CLI
credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html). For example, you could
set the credentials as the environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.

**Note #3**: Never hit `CTRL + C` or cancel a build once tests are running or the cleanup tasks won't run!

**Prerequisite**: The tests expect Terraform, Terragrunt, Packer, and/or Docker to already be installed and in your `PATH`.

To run all the tests:

```bash
go test -v -timeout 30m -p 1 ./...
```

To run the tests in a specific folder:

```bash
cd "<FOLDER_PATH>"
go test -timeout 30m
```

To run a specific test in a specific folder:

```bash
cd "<FOLDER_PATH>"
go test -timeout 30m -run "<TEST_NAME>"
```

### Versioning

This repo follows the principles of [Semantic Versioning](http://semver.org/). You can find each new release,
along with the changelog, in the [Releases Page](https://github.com/gruntwork-io/terratest/releases).

During initial development, the major version will be 0 (e.g., `0.x.y`), which indicates the code does not yet have a
stable API. Once we hit `1.0.0`, we will make every effort to maintain a backwards compatible API and use the MAJOR,
MINOR, and PATCH versions on each release to indicate any incompatibilities.
