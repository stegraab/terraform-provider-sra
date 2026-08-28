<!-- Modified by Stegra AB for the Stegra-maintained distribution. -->

# Contributing to the Stegra SRA Terraform Provider

This repository is a Stegra-maintained derivative of
[`BeyondTrust/terraform-provider-sra`](https://github.com/BeyondTrust/terraform-provider-sra).
Contributions should preserve upstream attribution and minimize unnecessary
differences so future upstream updates remain reviewable.

## Reporting issues

Report provider defects and feature requests in the
[Stegra repository](https://github.com/stegraab/terraform-provider-sra/issues).
Do not ask BeyondTrust Support to support Stegra-specific releases.

Security vulnerabilities in upstream BeyondTrust products or services should
follow [BeyondTrust's responsible disclosure
process](https://www.beyondtrust.com/security#disclosure). Vulnerabilities in
Stegra-specific provider changes should be reported through Stegra's approved
security channel rather than a public issue.

## Development

Terraform schemas should map directly to the relevant SRA Configuration API
contract. The appliance exposes its Configuration API documentation under
**Management > Security** in the `/login` interface for authorized
administrators.

Before submitting a pull request:

```shell
make generate
make unittest
```

Run the configured lint and static checks when available. New resource behavior
should have focused unit tests, including schema, mapping, import, and mocked
lifecycle behavior where practical.

## Live end-to-end tests

Tests under `./test` connect to an appliance and create and destroy SRA
resources. They are never part of ordinary pull-request CI. Run them only
through the protected manual workflow or locally against an explicitly
approved, non-production appliance:

```shell
make teste2e
```

## Upstream updates

The `upstream-main` branch mirrors BeyondTrust's main branch. Use
`scripts/sync_upstream_main.sh` to update it, then prepare an update on a
temporary `chore/sync-upstream-*` branch.

Use `git range-diff` as described in [UPSTREAM.md](UPSTREAM.md) to demonstrate
that the Stegra patch series, especially managed Group Policy support, remains
intact. Never publish upstream release tags as Stegra releases.

## Licensing

Contributions are accepted under the Apache License 2.0. Retain applicable
upstream copyright, patent, trademark, and attribution notices. Clearly mark
changes to upstream files and do not imply endorsement or support by
BeyondTrust Corporation.
