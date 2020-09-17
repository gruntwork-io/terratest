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

To help work around this problem, Terratest includes:

1. [`GetRecommendedInstanceType`](#getrecommendedinstancetype): A Go function that helps you pick a recommended instance type.
1. [`pick-instance-type`](#pick-instance-type): A CLI tool that helps you pick a recommended instance type.




## `GetRecommendedInstanceType`

`GetRecommendedInstanceType` takes in an AWS region and a list of EC2 instance types and returns the first instance 
type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
instance available in all AZs, this function exits with an error. 

Example usage:

```go
aws.GetRecommendedInstanceType(t, "eu-west-1", []string{"t2.micro", "t3.micro"})
// As of July, 2020, returns "t2.micro"

aws.GetRecommendedInstanceType(t, "ap-northeast-2", []string{"t2.micro", "t3.micro"})
// As of July, 2020, returns "t3.micro"
```   



## `pick-instance-type`

`pick-instance-type` is a CLI tool that you can download from the [Terratest releases 
page](https://github.com/gruntwork-io/terratest/releases) (click "Assets" under any release). It takes in an AWS 
region and a list of EC2 instance types and prints to `stdout` the first instance type in the list that is available in 
all Availability Zones (AZs) in the given region. If there's no instance available in all AZs, `pick-instance-type`
exits with an error.

Example usage:

```bash
# Data below is from July, 2020

$ pick-instance-type eu-west-1 t2.micro t3.micro
t2.micro

$ pick-instance-type ap-northeast-2 t2.micro t3.micro
t3.micro
```   
