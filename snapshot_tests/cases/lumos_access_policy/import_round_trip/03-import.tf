# The policy is found by its name and imported with the config it was created
# with, so the import plans no changes.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Policy Import"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

data "lumos_access_policies" "this" {
  name = "bb-${var.run_id} Imported"
}

import {
  to = lumos_access_policy.this
  id = one([for p in data.lumos_access_policies.this.access_policies : p.id if p.name == "bb-${var.run_id} Imported"])
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Imported"
  business_justification = "Snapshot test policy"
  status                 = "published"
  is_enabled             = true

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{
    id             = lumos_app.this.id
    is_preapproved = true
    entitlements   = [{ id = lumos_entitlement.this.id }]
  }]
}
