---
layout: collection-browser-doc
title: Behaviour changes
category: migrating-to-v2
excerpt: >-
  The two v2 changes the compiler will not find for you. Both are in the k8s module.
tags: ["migration", "v2"]
order: 403
nav_title: Documentation
nav_title_link: /docs/
---

Everything else in the v2 migration is a compile error. These two are not: your code builds and behaves differently.
Both are in `k8s`, so skip this page if you do not use it.

## Node addresses prefer `ExternalIP`

`k8s.FindNodeHostnameContextE`, and `GetServiceEndpoint` for a NodePort service, now return the `ExternalIP` recorded
on the Node object when one is present. They fall back to the internal hostname exactly as before when it is not.

Cloud controller managers record an instance's public IP as an `ExternalIP`. Terratest previously ignored that field
and, on AWS, called `ec2:DescribeInstances` to find the public IP instead. Reading the Node object is the same answer
without the API call, so:

- `ec2:DescribeInstances` is no longer needed for this path
- `k8s` no longer depends on the `aws` module, which removes 23 AWS service SDKs from its dependency graph

**What to check.** If your cluster advertises an `ExternalIP` and your test previously received an internal hostname,
it now receives the external address. That is the documented behaviour and almost certainly what you wanted, but it
is a different string. Tests that assert on the endpoint value, or that rely on reaching the node over its internal
address, are the ones to look at.

Signatures are unchanged. For clusters that do not advertise an `ExternalIP`, a new pair of functions takes
`*KubectlOptions` and consults a lookup you provide:

```go
options := k8s.NewKubectlOptions("", kubeconfig, "default")
options.NodePublicIPLookup = aws.GetPublicIpsOfEc2InstancesContextE

hostname, err := k8s.FindNodeHostnameWithOptionsContextE(t, ctx, options, node)
```

`NodePublicIPLookup` is only consulted after the Node's own `ExternalIP` has been checked, so most callers can leave
it nil.

## `KubectlOptions` carrying a `RestConfig` cannot be saved

`json.Marshal` of a `KubectlOptions` built by `NewKubectlOptionsWithRestConfig` now returns
`k8s.ErrRestConfigNotSerializable`. This affects `k8s.SaveKubectlOptions`, `teststructure.SaveTestData`, and any code
of your own that marshals options.

This already failed in v1, with an opaque `json: unsupported type: transport.WrapperFunc`, because `rest.Config`
holds func-typed fields that `encoding/json` rejects. What changed is that the error now says what is wrong and what
to do instead.

It is deliberately an error rather than silently dropping the config. Dropped, the reloaded options would carry no
cluster identity at all, fall back to the ambient kubeconfig, and run your test against a different cluster.

**What to do.** For staged tests, build options from a kubeconfig path or from in-cluster auth. Both round trip:

```go
options := k8s.NewKubectlOptions(contextName, configPath, namespace)
// or
options := k8s.NewKubectlOptionsWithInClusterAuth()
```

If you need a `rest.Config` at runtime, keep building one, but rebuild it in each stage rather than saving it.
