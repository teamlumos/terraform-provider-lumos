# requires: GROUP_ID
# The settings are imported by the entitlement ID with the config they were
# created with, so the import plans no changes.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Import"
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
  entitlement_id                         = lumos_entitlement.this.id
  visible_in_appstore                    = false
  allowed_requester_group_uuids_override = [var.group_id]

  depends_on = [lumos_app_appstore_settings.this]
}

data "lumos_current_user" "me" {}
