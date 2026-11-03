<!-- Copyright 2026 Stegra AB. SPDX-License-Identifier: Apache-2.0 -->

# Upstream provenance and maintenance

This Stegra-maintained derivative of
[`BeyondTrust/terraform-provider-sra`](https://github.com/BeyondTrust/terraform-provider-sra)
retains the upstream Apache License 2.0 and applicable attribution. It is not
affiliated with, endorsed by, or supported by BeyondTrust Corporation.

## Current base and patches

- Upstream release: `v1.5.0` (`afe77d223585790d088e6af8ef816496df24d578`).
- Terraform Registry address: `registry.terraform.io/stegraab/sra`.
- Group Policy, Group Policy Member, ephemeral Vault password support, and
  read-only Vault update handling use the upstream implementations and tests.
- Remaining functional patches: `sra_security_provider` lookup by name and
  plan modifiers preserving computed Vault account values on unrelated updates.
- Distribution patches retain Stegra branding, signing, release versioning,
  code ownership, and manual appliance-test workflows.

## Maintenance

`upstream-main` mirrors BeyondTrust's selected upstream main commit; `main`
is the Stegra release branch. Prepare updates on `chore/sync-upstream-*`
branches and compare the remaining diff against the new upstream base.
Use `scripts/sync_upstream_main.sh` to update the mirror.

Before replacing an overlapping patch, compare schema, import, lifecycle,
and state compatibility. Retain upstream resource registration and lifecycle
tests alongside tests for remaining Stegra behavior.

## Releases

`UPSTREAM_VERSION` records the base version. The manual release workflow
uses Stegra's signing key and creates the next `vX.Y.Z-stegra.N` tag.
Only Stegra release tags may be pushed to `origin`.

Live end-to-end tests create and destroy appliance resources. They are manual,
use protected GitHub environments, and must never target production.
