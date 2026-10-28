# requires: GROUP_ID
# The settings are imported with a config that leaves out the group override,
# so the import plans to clear it and inherit the app's groups.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Import Clears"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = true
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

import {
  to = lumos_entitlement_appstore_settings.this
  id = lumos_entitlement.this.id
}

resource "lumos_entitlement_appstore_settings" "this" {
  entitlement_id      = lumos_entitlement.this.id
  visible_in_appstore = false

  depends_on = [lumos_app_appstore_settings.this]
}

data "lumos_current_user" "me" {}
