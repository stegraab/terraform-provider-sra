# Modified by Stegra AB for the Stegra-maintained distribution.
# SPDX-License-Identifier: Apache-2.0
resource "sra_group_policy" "example" {
  name = "Terraform Managed Group Policy"

  perm_collaborate          = true
  perm_session_idle_timeout = 900
}
