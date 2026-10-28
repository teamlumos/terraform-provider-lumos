# The settings are imported by the app ID with the config they were created
# with, so the import plans no changes.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Import Round Trip"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

import {
  to = lumos_app_appstore_settings.this
  id = lumos_app.this.id
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "search-only"
  allow_permission_selection          = false
  allow_multiple_permission_selection = false
  request_instructions                = "Snapshot test instructions"
}

data "lumos_current_user" "me" {}
