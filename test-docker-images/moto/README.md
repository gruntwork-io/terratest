# Moto

This docker image runs [moto](https://github.com/spulec/moto) as a service. We will use Moto as a local service that
accepts AWS API calls and returns valid API responses. Moto can be used with any AWS SDK, including the [awscli](
https://aws.amazon.com/cli/)[].

This Docker image is expected to run alongside a cluster of docker containers to represent a "local AWS".

### About Moto

Moto was originally written as a way to mock out the Python [boto](https://github.com/boto/boto3) library, the official
AWS SDK for Python. It was motivated by the need to write automated tests for boto. But since the AWS API cannot run
"locally", Moto was written to mock the responses of the AWS API.

Moto works by receiving AWS API requests across a wide variety of AWS services, including most of EC2. It will then store
the requested AWS resource in memory and allow you to query that AWS resource using standard AWS API calls. There is
no actual VM created, or other actual resource created.

### Motivation

As part of writing [Unit Tests with Terratest](/README.md#unit-tests), we need a way to run our services in a Docker
container. But this presents a new challenge: Almost all our cluster-based setups query the AWS APIs to obtain metadata
about the EC2 Instance on which they're running. How can we simulate these API calls in a local environment? Moto seems
to meet this use case perfectly.

## Usage

### Building and Pushing a New Docker Image to Docker Hub

This Docker image should publicly accessible via Docker Hub at https://hub.docker.com/r/gruntwork/moto/. To build and
upload it:

1. `docker build -t gruntwork/moto:v1 .`
1. `docker push gruntwork/moto:v1`

#### Run a Docker container

```
docker run -p 5000:5000 gruntwork/moto moto_server ec2 --host 0.0.0.0
```

This runs the `moto` service as a RESTful API, specially for the AWS EC2 API with support for accepting connections from
any IP address (versus just from localhost). For additional information:
- See the [moto stand-alone server usage docs](https://github.com/spulec/moto#stand-alone-server-mode)
- See [which AWS services are supported](https://github.com/spulec/moto#in-a-nutshell)

#### Make AWS API calls against Moto

Because Moto exposes an API that is intended to be identical to the official AWS API, you can use any any AWS SDK against
it, including the AWS CLI, AWS SDK for Go, `curl` calls, or any other AWS API library. The only difference is that you
must explicitly set the "endpoint uRL" to point to the Moto server instead of the official AWS API. Changing this setting
will be different for each AWS SDK, but for the AWSCLI, you can simplify specify the `--endpoint-url` argument as follows:

```
aws --region "us-west-2" --endpoint-url="http://localhost:5000" ec2 run-instances --image-id ami-abc12345 --tag-specifications 'ResourceType=instance,Tags=[{Key=ServerGroupName,Value=josh}]'
```

Note that Moto supports all AWS regions, and will automatically create a VPC with default subnets for you!

## Other Solutions Considered

### LocalStack

[LocalStack](https://localstack.cloud/) is a (coming soon) commercial service intending to offer "local AWS as a service".
It is based on the open source [localstack](https://github.com/localstack/localstack) project.

Local stack seems to offer a [small set of advantages](https://github.com/localstack/localstack#why-localstack) over
Moto, including throwing periodic errors to simulate a real-world cloud environment. But Localstack doesn't implement
100% of the APIs implemented by Moto (including EC2, the most important one for us!), its docker image is ~500 MB, it
doesn't appear to be an active commercial entity, and Moto supports a RESTful API already.

For these reasons, Moto is the better fit for our needs.