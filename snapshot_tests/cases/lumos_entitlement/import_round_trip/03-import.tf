# The entitlement is found by its label and imported with the config it was
# created with, so the import plans no changes.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Import"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_app_entitlements" "this" {
  app_id = lumos_app.this.id
  search = "bb-${var.run_id} Imported"
}

import {
  to = lumos_entitlement.this
  id = one([for e in data.lumos_app_entitlements.this.entitlements : e.id if e.label == "bb-${var.run_id} Imported"])
}

resource "lumos_entitlement" "this" {
  app_id                           = lumos_app.this.id
  label                            = "bb-${var.run_id} Imported"
  description_override             = "Imported entitlement."
  custom_approval_message_override = "Your imported access is ready."
}
