---
layout: collection-browser-doc
title: Migrating Azure tests to v1.0.0
category: migrating-to-v1
excerpt: How to adapt your Azure tests to Terratest v1.0.0's new module surface.
tags: ["v1", "azure", "migration"]
order: 300
nav_title: Documentation
nav_title_link: /docs/
---

The `modules/azure` package received the largest set of breaking changes in
the v1.0.0 release. This guide walks through what changed and how to update
your tests.

## Why we migrated

The previous version of `modules/azure` was built on
`github.com/Azure/azure-sdk-for-go/services/...`, the legacy "track 1" Azure
SDK. Microsoft has archived that SDK; it no longer receives feature updates,
bug fixes, or security patches.

v1.0.0 moves the entire package to
`github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/...`, the modular,
actively maintained "track 2" SDK. The new SDK has a different shape: typed
resource clients are produced by a `ClientFactory`, response payloads put
most fields on a nested `Properties` struct, and pagination uses pagers
instead of iterators. We took the opportunity to land a few small API
cleanups at the same time so v1.0.0 ships a coherent, stable surface.

## What changed at a glance

- All Azure service code now imports `sdk/resourcemanager/<service>/arm<service>` packages instead of `services/<service>/mgmt/<api-version>/<service>`.
- Resource fields moved under `.Properties` (e.g. `vm.StorageProfile` is now `vm.Properties.StorageProfile`).
- Iterator-based list calls (`NextWithContext`) are replaced with pagers (`NewListPager` / `More` / `NextPage`).
- 8 deprecated `Get*ClientE` client-getter functions were removed; the `Create*ClientE` replacements have been around for a while.
- 4 `CreateNew*ClientContextE` factories were renamed to `Create*ClientContextE`. The old names remain as deprecated aliases.
- `NsgRuleSummary.SourceAdresssPrefixes` (triple-s typo) was renamed to `SourceAddressPrefixes`.
- A new `*WithClient` family of functions was added so tests can inject a fake or pre-built SDK client (useful with the Azure SDK's `azfake` package).
- `GetVirtualMachineImage` / `GetVirtualMachineImageE` now return `*VMImage` instead of `VMImage`.

## Updating SDK imports

Most Terratest users do not import the underlying Azure SDK directly,
because Terratest wraps it. If you only call `terratest/modules/azure`
helpers, you can usually skip this section. If your tests do import the SDK
(for example to construct request objects or assert on returned types),
update imports as follows:

| Old (`services/...`) | New (`sdk/resourcemanager/...`) |
| --- | --- |
| `github.com/Azure/azure-sdk-for-go/services/compute/mgmt/.../compute` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6` |
| `github.com/Azure/azure-sdk-for-go/services/network/mgmt/.../network` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6` |
| `github.com/Azure/azure-sdk-for-go/services/storage/mgmt/.../storage` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage` |
| `github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/.../documentdb` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3` |
| `github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/.../servicebus` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2` |
| `github.com/Azure/azure-sdk-for-go/services/preview/containerservice/.../containerservice` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6` |
| `github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/.../containerregistry` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry` |
| `github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/.../containerinstance` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2` |
| `github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/.../operationalinsights` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2` |
| `github.com/Azure/azure-sdk-for-go/services/resources/mgmt/.../subscriptions` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions` |
| `github.com/Azure/azure-sdk-for-go/services/resources/mgmt/.../resources` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources` |
| `github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/.../privatedns` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns` |
| `github.com/Azure/azure-sdk-for-go/profiles/latest/frontdoor/mgmt/frontdoor` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor` |
| `github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor` |
| `github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/.../recoveryservices` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices` |
| `github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/.../backup` | `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v4` |

Once imports are updated, expect three follow-on edits per file:

1. **Type names** lose their old prefix and gain `arm`. For example
   `compute.VirtualMachine` becomes `armcompute.VirtualMachine`,
   `network.SecurityGroup` becomes `armnetwork.SecurityGroup`, and
   `storage.Account` becomes `armstorage.Account`.
2. **Field access** moves under `.Properties`. For example
   `vm.StorageProfile` becomes `vm.Properties.StorageProfile`, and
   `registry.LoginServer` becomes `registry.Properties.LoginServer`.
3. **List iteration** moves to the pager pattern: replace
   `iterator.NextWithContext(ctx)` loops with
   `for pager.More() { page, err := pager.NextPage(ctx); ... }`.

## Renamed factory functions

Four client factories were renamed to drop the redundant `New` (a
`Create*New*Client` reads as redundant). The old names remain as deprecated
aliases for one minor release; please update at your convenience.

| Old name | New name |
| --- | --- |
| `CreateNewNetworkInterfacesClientE` | `CreateNetworkInterfacesClientE` |
| `CreateNewNetworkInterfacesClientContextE` | `CreateNetworkInterfacesClientContextE` |
| `CreateNewNetworkInterfaceIPConfigurationClientE` | `CreateNetworkInterfaceIPConfigurationClientE` |
| `CreateNewNetworkInterfaceIPConfigurationClientContextE` | `CreateNetworkInterfaceIPConfigurationClientContextE` |
| `CreateNewSubnetClientE` | `CreateSubnetClientE` |
| `CreateNewSubnetClientContextE` | `CreateSubnetClientContextE` |
| `CreateNewVirtualNetworkClientE` | `CreateVirtualNetworkClientE` |
| `CreateNewVirtualNetworkClientContextE` | `CreateVirtualNetworkClientContextE` |

## Removed deprecated functions

The previous release marked eight client-getter functions for removal
("`TODO: remove in next version`"). v1.0.0 is that version. Each removed
function has a long-standing `Create*ClientE` replacement.

| Removed | Replacement |
| --- | --- |
| `GetAvailabilitySetClientE` | `CreateAvailabilitySetClientE` |
| `GetDiskClientE` | `CreateDisksClientE` |
| `GetDiagnosticsSettingsClientE` | `CreateDiagnosticsSettingsClientE` |
| `GetVMInsightsClientE` | `CreateVMInsightsClientE` |
| `GetActivityLogAlertsClientE` | `CreateActivityLogAlertsClientE` |
| `GetResourceGroupClientE` | `CreateResourceGroupClientE` |
| `GetStorageAccountClientE` | `CreateStorageAccountClientE` |
| `GetStorageBlobContainerClientE` | `CreateStorageBlobContainerClientE` |

The replacements take the same arguments and return the same client type
(now from the new SDK). The rename is mechanical: `Get` → `Create`. Note
that `GetDiskClientE` becomes `CreateDisksClientE` (plural) to match the
underlying SDK's `DisksClient` type.

## Typo fix on `NsgRuleSummary`

`NsgRuleSummary.SourceAdresssPrefixes` (note the three s's) was renamed to
the correctly-spelled `SourceAddressPrefixes`. The field type
(`[]string`) is unchanged. Update any code that read or set this field:

```go
// Before
for _, prefix := range rule.SourceAdresssPrefixes {
    // ...
}

// After
for _, prefix := range rule.SourceAddressPrefixes {
    // ...
}
```

The paired `DestinationAddressPrefixes` field was already spelled
correctly and is unchanged.

## `VMImage` is now a pointer

`GetVirtualMachineImage` and `GetVirtualMachineImageE` now return
`*VMImage` instead of `VMImage`, matching every other resource getter in
the package.

```go
// Before
img := azure.GetVirtualMachineImage(t, vmName, rg, sub)
fmt.Println(img.Publisher)

// After
img := azure.GetVirtualMachineImage(t, vmName, rg, sub)
if img != nil {
    fmt.Println(img.Publisher)
}
```

If the resource cannot be loaded the function still fails the test, so
the `nil` guard is precautionary.

## New `WithClient` variants for testability

v1.0.0 adds a parallel family of `*WithClient` functions across all Azure
modules. Each one accepts a pre-built SDK client and a
`context.Context`, so you can drive Terratest helpers against the Azure
SDK's `azfake` fake-server framework in unit tests:

```go
import (
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6/fake"
)

// Build a fake client that returns a canned response.
fakeServer := fake.DisksServer{ /* ... */ }
client, _ := armcompute.NewDisksClient("sub-id", nil, &arm.ClientOptions{
    ClientOptions: azcore.ClientOptions{
        Transport: fake.NewDisksServerTransport(&fakeServer),
    },
})

disk, err := azure.GetDiskWithClient(ctx, client, resourceGroup, diskName)
```

This is purely additive: the existing `*ContextE` functions still work and
build their own clients from ambient credentials. Use `WithClient`
variants only if you need test injection.

## Search-and-replace cheatsheet

Most projects can do the bulk of the migration with a few find/replace
passes. The snippets below cover the most common edits:

```bash
# Removed Get*ClientE functions -> Create*ClientE
sd 'GetAvailabilitySetClientE\b'      'CreateAvailabilitySetClientE'      $(rg -l 'GetAvailabilitySetClientE')
sd 'GetDiskClientE\b'                 'CreateDisksClientE'                $(rg -l 'GetDiskClientE')
sd 'GetDiagnosticsSettingsClientE\b'  'CreateDiagnosticsSettingsClientE'  $(rg -l 'GetDiagnosticsSettingsClientE')
sd 'GetVMInsightsClientE\b'           'CreateVMInsightsClientE'           $(rg -l 'GetVMInsightsClientE')
sd 'GetActivityLogAlertsClientE\b'    'CreateActivityLogAlertsClientE'    $(rg -l 'GetActivityLogAlertsClientE')
sd 'GetResourceGroupClientE\b'        'CreateResourceGroupClientE'        $(rg -l 'GetResourceGroupClientE')
sd 'GetStorageAccountClientE\b'       'CreateStorageAccountClientE'       $(rg -l 'GetStorageAccountClientE')
sd 'GetStorageBlobContainerClientE\b' 'CreateStorageBlobContainerClientE' $(rg -l 'GetStorageBlobContainerClientE')

# CreateNew*Client renames
sd 'CreateNewNetworkInterfacesClient'              'CreateNetworkInterfacesClient'              $(rg -l 'CreateNewNetworkInterfacesClient')
sd 'CreateNewNetworkInterfaceIPConfigurationClient' 'CreateNetworkInterfaceIPConfigurationClient' $(rg -l 'CreateNewNetworkInterfaceIPConfigurationClient')
sd 'CreateNewSubnetClient'                         'CreateSubnetClient'                         $(rg -l 'CreateNewSubnetClient')
sd 'CreateNewVirtualNetworkClient'                 'CreateVirtualNetworkClient'                 $(rg -l 'CreateNewVirtualNetworkClient')

# NsgRuleSummary typo
sd 'SourceAdresssPrefixes' 'SourceAddressPrefixes' $(rg -l 'SourceAdresssPrefixes')
```

For the SDK migration itself we recommend doing one Azure service at a
time, starting with the import path, then fixing the resulting compile
errors (type names, `.Properties` access, pager loops). The Go compiler
is the most reliable migration tool here.

## Need help

Open an issue on the Terratest repo with the `azure` label and a snippet
of the failing code. If you spot a gap in this guide, please send a PR
against `docs/_docs/03_migrating-to-v1/azure.md`.
