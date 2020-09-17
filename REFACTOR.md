# Terratest refactor

Terratest started out as a set of Bash scripts we were using at Gruntwork to test some of our Terraform code. As the
amount of Terraform code grew, it was getting trickier and trickier to test it with Bash, so we rewrote those scripts
in Go. Over time, this Go code grew into a library called Terratest, which contains collection of utilities that we
use to test all aspects of our [Infrastructure as Code Library](https://gruntwork.io/infrastructure-as-code-library/).

We developed patterns to test Terraform configurations, Packer templates, Docker images, SSH access, AWS APIs, shell
commands, and much more. We built this library because we couldn't find any existing tools out there that could do
the type of real-world testing we needed. It turns out many other companies want to do this type of testing too, so
now it's time to open source Terratest.

This library grew organically, so it needs lots of refactoring, cleanup, and documentation to be useful to people
outside of Gruntwork. This document lays out the refactoring we are planning to (a) get feedback and (b) document what
changed so that when we update our code, we know how to deal with the backwards incompatibilities.




## Name change

"Terratest" made sense as the name for this library when it was all about testing Terraform code, but now this library
also can help you test Packer templates, Docker images, and much more. I propose that we rename it. Some ideas:

- grunt-test
- test-grunt
- gruntUnit
- iac-test
- infratest




## Build updates

1. Move from Glide to Dep
1. Move from CircleCI 1.0 to 2.0
1. Add support for Google Cloud Platform
1. Move from Dep to Go Modules.


## Folder structure

Change to the same folder structure we use for just about all other Gruntwork repos:

- `examples`: This will contain a number of real-world examples of code you might want to test with Terratest, such as
  Terraform modules, Packer templates, and Docker images. The `test` folder (described below) shows how to use
  Terratest to test these examples.

- `modules`: The Terratest source code. Move all the `.go` files and packages into this folder so it's easier to browse
  the repo. That does mean all Terratest imports will have to be updated to
  `github.com/gruntwork-io/terratest/modules/xxx`. Unit tests for the Go source code will be in this folder too (e.g., the
  unit test for `foo.go` will be in `foo_test.go`).

- `test`: This will contain the automated tests for the examples in the `examples` folder. These will act both as an
  example of how to use Terratest, as well as integration tests for the library.




## Documentation

Update the root `README.md` with documentation that shows how to use Terratest:

- Overview of what Terratest is.
- Link to blog post we'll write about Terratest (this blog post is a TODO for after this refactor).
- Discuss some of the challenges of testing infrastructure code: i.e., lack of "localhost," lock of "unit tests,"
  slowness, brittleness.
- Discuss the value of doing this testing despite the challenges: i.e., there is no way to maintain lots of
  infrastructure code without tests, building reusable, tested, versioned modules changes how you manage
  infrastructure.
- Discussion of test strategies: using Docker for local testing, test stages, retries, mocks, small modules, test
  pyramid, cleanup, `cloud-nuke`.
- Point to `examples` folder for real-world code you may want to test and `test` folder for examples of how to use
  Terratest to test that code.
- Overview of Terratest packages. Explain what each top-level package in Terratest does. We can't do a method-by-method
  breakdown, as that would go out of date immediately, so instead, link to the appropriate `examples` subfolder that
  shows real-world usage of that package.
- In the future: links to our open source repos (Vault, Consul, Nomad, Couchbase, etc) that show how we use Terratest
  with our own code. We can't add this until we update those open source repos to this refactored version of Terratest
  so the code matches up.




## Package-by-package refactor

I've gone through each of the packages in Terratest and took down some notes on cleanup we need to do. This is not a
comprehensive list, as things will become clearer once I actually start doing the work.

In fact, my plan is to first create all the examples in the `examples` folder, then write tests for them in the `test`
folder using "wishful thinking" (in the
[SICP](https://www.amazon.com/Structure-Interpretation-Computer-Programs-Engineering/dp/0262510871) sense), where I
come up with the test API I want to have for doing the testing, and then go and refactor the Terratest code to match.


### Root package

We have a lot of stuff in the root package and I propose moving all of it out into appropriate sub-packages:

- `apply.go`, `apply_and_destroy.go`, `destroy.go`, `output.go`, and `output_test.go` will all be moved into
  `modules/terraform`, as they are all specific to testing Terraform code.

- I propose deleting `rand_resources.go` and `rand_resource_test.go` and extracting its logic into other places. The
  `RandomResourceCollection` ended up being a, well, random collection of resources, most of which don't apply to most
  of our tests, and certainly won't apply to tests written by the open source community. Here's what
  `RandomResourceCollection` contains and what I propose to do with it:

    - `UniqueId`: We already have a separate method for generating a unique ID and we can pass it around as a `string`.

    - `AwsRegion`: This is only needed for AWS tests. We want to expand Terratest to support other clouds, so it needs
      to be separated anyway. Code that needs an AWS region should call a method in the `modules/aws` package to pick a
      random AWS region (passing in a list of forbidden regions, if necessary) and can pass that around as a `string`.

    - `KeyPair`: This is only needed for a small percentage of our AWS tests that deploy EC2 Instances and SSH to them.
      Those tests should call a standalone method in the `modules/aws` package to generate this `KeyPair` when they need it,
      instead of us assuming every single test needs it.

    - `AmiId`: We used to look up vanilla Ubuntu or Amazon Linux AMI IDs and put them in this field, but now that
      Terraform has `data` sources and Packer has `source_ami_filter`, this is no longer necessary. We can keep the
      methods around to find Ubuntu or Amazon Linux AMI IDs for tests that need them, but there's no need to assume
      every single test needs this.

    - `AccountId`: Our Terraform examples used to require an account ID to be passed in. We now avoid this to make the
      examples easier to use, and fetch it automatically using Terraform's `aws_caller_identity` data source if it's
      absolutely necessary. Code that needs an account ID should call a method in the `modules/aws` package to fetch it,
      but we shouldn't assume every single test needs it.

    - `SnsTopicArn`: A very, very small percentage of our tests needed an SNS topic passed in. Those tests should call
      a method in the `modules/aws` package to create this topic instead of us assuming every single test needs it.

- I propose moving `terratest_options.go` to `modules/terraform/options.go` and renaming the struct within it from
  `TerratestOptions` to `Options`, since this is solely used for testing Terraform code. We should also rename
  `TemplatePath` to `TerraformDir`, as `.tf` files are technically called "configurations" and not "templates".

- `url_checker.go` will be deleted. It's too hard-coded for one specific type of check. The reuse value is limited and
  it's not obvious the code exists, so it's best for the test cases to reimplement this themselves, with their specific
  needs, even if it's a tiny bit less DRY.


### _docker-images package

- Rename to `test-docker-images` to make it clearer these are only used for testing.
- Use these Docker images in the `examples` folder to show how to do "unit tests" for Packer templates.
- Follow-up PR: build and push a new version of these Docker images on each release?
- Follow-up PR: tag each new Docker image with a unique version number (e.g., sha1 of commit).


### `aws` package

- Right now, much of this code has no unit tests, since it relies on resources in AWS. By adding an `examples` folder
  that deploys real resources in AWS, we will be able to test this code better, _and_ show users how to use this code!

- `ami.go`: Update these methods to use the AWS APIs to find the latest Ubuntu / Amazon Linux AMI IDs instead of
  hard-coding them.

- `kms.go`: What to do about `GetDedicatedTestKeyArn`? For tests that use KMS, we don't want to create a new CMK each
  time the test runs, as AWS charges $1/month for CMKs, even if you delete them immediately after use. This method
  currently assumes we have a key called `alias/dedicated-test-key` in every AWS region. Should we leave it as-is and
  document it for Terratest users that want to follow a similar pattern? Or perhaps read the key name from an env var?

- `region.go`: What should we do about `GetGloballyForbiddenRegions`? Right now, it's hard-coded to include `us-west-2`
  as a globally forbidden region, as Josh is running his personal blog there. Obviously, we don't want that in the open
  source version. Josh, can you finally migrate your blog out of there so we don't have to have this exception?


### `log` package

- Rename to `logger`. That way, we don't have to alias it as `terralog` all over our test code.
- Change what the package does. Instead of creating a custom `*log.Logger` and passing it around, we are going to have
  a `Log` and `Logf` method you can call from anywhere. To use those methods, you have to pass them a `*testing.T`,
  which they will use to read out the test name. We already pass `*testing.T` to almost all of our test methods, so
  this reduces the number of arguments by one.


### `parallel` package

- I propose removing this package entirely. Now that go has [subtests](https://blog.golang.org/subtests) that you can
  easily run with `t.Run()` and parallelize with `t.Parallel()`, I think that's a cleaner way of handling parallelism
  than this custom package.


### `packer` package

- Rename `PackerOptions` to `Options` (the package name is already `packer`).


### `resources` package

- `base_resources.go` is no longer necessary if we remove `RandomResourceCollection`.
- `exclusions.go` is not used much and very out of date.
- `terraform_options.go` is hard-coded to how we do things at Gruntwork, but won't apply to many other users.


### `terraform` package

- `apply.go`: Remove `terraformDebugEnv` and instead make it easy to pass a map of env vars to the `Apply` method.
  Refactor `ApplyAndGetOutputWithRetry` to accept a list of errors on which to retry and how many retries to do.


### `test-util` package

- `dummy_server.go`: Move into the `modules/http` package.
- Remove `test-util` since that would leave it empty!


### `util` package

- `collections.go`: Move into its own `modules/collections` package.
- `keygen.go`: Move into `modules/ssh` package.
- `network.go`: Move into `modules/aws` package.
- `sleep.go`: Remove. Didn't even know we had this and doubt it gets much use!
- `random.go`: Move into its own `modules/random` package.
- `retry.go`: Move into its own `modules/retry` package.




## Error handling

I am updating most of the methods to support handling errors in one of two ways:

1. Each method `foo` will take in a `*testing.T` and upon hitting an error, call `t.Fatal`.
1. Each method `fooE` will explicitly return any errors it hits and NOT call `t.Fatal`.

Example:

```go
func GetCurrentBranchName(t *testing.T) string {
	out, err := GetCurrentBranchNameE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func GetCurrentBranchNameE(t *testing.T) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}
```

In most places in our code, we will use `GetCurrentBranchName`, which will call `t.Fatal` if it hits any errors. This
is typically the behavior we want anyway, and not having to deal with a returned error will keep our code smaller and
easier to read. However, in those cases where we may want to get the original error back and not fail the test
immediately, we can use `GetCurrentBranchNameE`.




## Other thoughts on the refactor

- Updating to the refactored version of Terratest will be a pain that requires lots of search & replace. But in the
  long term, it seems like worthwhile cleanup.

- There are a bunch of patterns we often end up using throughout our tests that would be good to copy into Terratest.
  Anyone remember what those are off the top of our head?

- Having two copies of each method (`foo` and `fooE`) is a bit tedious, but the `foo` variety is essentially the same
  boilerplate everywhere, so it only increases the maintenance burden on Terratest library maintainers a little, but
  it improves code readability for all Terratest users enormously.
