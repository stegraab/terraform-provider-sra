<!-- Copyright 2026 Stegra AB. SPDX-License-Identifier: Apache-2.0 -->

# Upstream provenance and maintenance

This repository is a Stegra-maintained derivative of
[`BeyondTrust/terraform-provider-sra`](https://github.com/BeyondTrust/terraform-provider-sra).
It retains the upstream Apache License 2.0 and applicable attribution.

## Initial Stegra release line

- Upstream base commit: `951cbbd605ade1bad782c8770b385c28084bbb98`
- Nearest upstream release: `v1.3.0` (`72` upstream commits before the
  selected base)
- Managed group policy commit: `9a12478260900142bb3401fe9c71b42e0417d806`
- Group policy member and PRA validation commit: `40bfc3584e961fe856f5050931d34ec41cb77590`
- Terraform Registry address: `registry.terraform.io/stegraab/sra`

The provider is maintained and supported by Stegra AB. It is not affiliated
with, endorsed by, or supported by BeyondTrust Corporation.

## Branches

- `upstream-main` is an exact mirror of the selected BeyondTrust upstream main
  branch. Do not commit to it.
- `main` is the Stegra release branch.
- Upstream updates are prepared on `chore/sync-upstream-*` branches and merged
  only after review.

Run `scripts/sync_upstream_main.sh` from a clean working tree to update the
mirror branch. The script fetches upstream tags for local comparison but does
not publish them to the Stegra repository. Only Stegra release tags may be
pushed to `origin`.

For every upstream update, compare the old and new Stegra patch series before
merging:

```shell
git range-diff \
  <old-upstream>..<old-stegra-main> \
  <new-upstream>..<sync-candidate>
```

The resource registration and lifecycle tests for `sra_group_policy` and
`sra_group_policy_member` are release contracts. They must remain enabled even
if upstream later adds overlapping functionality. Replace Stegra-specific code
only after schema, import, lifecycle, and state-compatibility behavior has been
compared.

## Release policy

Releases use Stegra's signing key and append `-stegra.N` to the selected
upstream version, for example `v1.3.0-stegra.1`. `UPSTREAM_VERSION` records the
base version. The manual release workflow increments `N` and creates the tag;
unsuffixed upstream tags never trigger Stegra releases.

Live end-to-end tests create and destroy appliance resources. They are manual,
must use protected GitHub environments, and must never target a production
appliance.
