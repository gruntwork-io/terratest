---
layout: collection-browser-doc
title: Picking EC2 instance types
category: testing-best-practices
excerpt: >-
  Pick EC2 instance types that are available in the current AWS region.
tags: ["testing-best-practices", "aws", "ec2"]
order: 213
nav_title: Documentation
nav_title_link: /docs/
---

It's common to want to test infrastructure code that deploys [EC2 instances](https://aws.amazon.com/ec2/) into AWS.
There are many different [instance types](https://aws.amazon.com/ec2/instance-types/), but not all instance types
are available in all [regions or availability zones
(AZs)](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html). For example,
`t3.micro` is sometimes available only in newer AZs, while `t2.micro` is sometimes only available in older AZs. If you
are testing code that needs to deploy a "small" instance across many regions, this can make it tricky to know which
region to pick.

To help work around this problem, Terratest includes
[`GetRecommendedInstanceTypeContext`](#getrecommendedinstancetypecontext), a Go function that helps you pick a
recommended instance type.

## `GetRecommendedInstanceTypeContext`

`GetRecommendedInstanceTypeContext` takes in an AWS region and a list of EC2 instance types and returns the first
instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
instance available in all AZs, this function exits with an error.

Example usage:

```go
ctx := t.Context()

aws.GetRecommendedInstanceTypeContext(t, ctx, "eu-west-1", []string{"t2.micro", "t3.micro"})
aws.GetRecommendedInstanceTypeContext(t, ctx, "ap-northeast-2", []string{"t2.micro", "t3.micro"})
```
