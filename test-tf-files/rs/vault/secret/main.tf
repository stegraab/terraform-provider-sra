# Modified by Stegra AB for the Stegra-maintained distribution.
# SPDX-License-Identifier: Apache-2.0
terraform {
  # This module is now only being tested with Terraform 1.1.x. However, to make upgrading easier, we are setting 1.0.0 as the minimum version.
  required_version = ">= 1.0.0"

  required_providers {
    sra = {
      source = "stegraab/sra"
    }
  }
}

module "account" {
  source      = "../user_pass_account"
  random_bits = var.random_bits
  name        = var.name
}

data "sra_vault_secret" "secret" {
  id = module.account.item.id
}

module "ssh" {
  source      = "../ssh_account"
  random_bits = "${var.random_bits}2"
  name        = "${var.name}2"
}

data "sra_vault_secret" "secret_ssh" {
  id = module.ssh.item.id
}

data "sra_vault_secret" "secret_ssh_ca" {
  id = module.ssh.stand_alone_ca.id
}
