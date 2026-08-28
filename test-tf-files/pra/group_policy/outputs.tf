# Modified by Stegra AB for the Stegra-maintained distribution.
# SPDX-License-Identifier: Apache-2.0
output "item" {
  value = sra_group_policy.item
}

output "listed" {
  value = data.sra_group_policy_list.listed.items
}
