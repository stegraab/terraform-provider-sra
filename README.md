<!-- Modified by Stegra AB for the Stegra-maintained distribution. -->

# Stegra SRA Terraform Provider

The Stegra SRA Provider allows [Terraform](https://terraform.io) to manage
resources in BeyondTrust [Secure Remote
Access](https://www.beyondtrust.com/secure-remote-access) products through the
Configuration API. It supports Remote Support and Privileged Remote Access,
subject to the API capabilities exposed by the appliance.

This is a Stegra-maintained derivative of the
[BeyondTrust SRA Terraform Provider](https://github.com/BeyondTrust/terraform-provider-sra).
It is not affiliated with, endorsed by, or supported by BeyondTrust
Corporation. Issues with this distribution should be reported in the
[Stegra repository](https://github.com/stegraab/terraform-provider-sra/issues),
not to BeyondTrust Support.

## Why this distribution exists

Stegra maintains this provider so required infrastructure capabilities can be
released on a predictable schedule while useful changes continue to be offered
upstream. The initial Stegra release adds managed Group Policy and Group Policy
Member resources while preserving the upstream provider's resources and data
sources.

The upstream base and Stegra patch history are recorded in
[UPSTREAM.md](UPSTREAM.md).

## Usage

```terraform
terraform {
  required_providers {
    sra = {
      source  = "stegraab/sra"
      version = "~> 1.4"
    }
  }
}

provider "sra" {}
```

The provider reads API connection settings from these environment variables:

- `BT_API_HOST`
- `BT_CLIENT_ID`
- `BT_CLIENT_SECRET`

They may also be configured explicitly:

```terraform
provider "sra" {
  host          = "example.beyondtrustcloud.com"
  client_id     = var.bt_client_id
  client_secret = var.bt_client_secret
}
```

The API account requires **Allow Access** permission for the Configuration API.
Vault resources additionally require **Manage Vault Accounts** permission.

## Compatibility

The upstream provider requires Remote Support or Privileged Remote Access
23.2.1 or later. Earlier appliance versions may return unsupported fields or
errors.

Managed Group Policy support is based on the bundled PRA OpenAPI v1.10
contract. Live PRA validation has covered import, no-change planning, create,
update, and delete. Remote Support uses the shared Group Policy model but has
not received equivalent live acceptance coverage in this distribution.

Some appliance versions return fields that are absent from the bundled OpenAPI
schema. For example, one validated PRA appliance returned
`perm_edit_group_policy_memberships`. Unknown response fields are ignored by
the client, but new fields cannot be configured until they are represented in
the bundled schema and provider model.

The Group Policy resource rejects enabled Jump permissions when
`perm_access_allowed` is explicitly `false`, matching PRA normalization and
preventing inconsistent post-apply state.

## Development

```shell
make generate
make unittest
```

The tests under `./test` are live end-to-end tests. They create and destroy
appliance resources and must only be run manually against an approved,
non-production appliance.

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution and upstream-sync
guidance.

## License and attribution

This distribution is licensed under the Apache License 2.0. See
[LICENSE.md](LICENSE.md) and [NOTICE](NOTICE). Modified files and Stegra
additions are identified by Git history and the provenance recorded in
[UPSTREAM.md](UPSTREAM.md).
