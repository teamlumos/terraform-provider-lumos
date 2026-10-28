# The settings are imported with a config that leaves out
# request_instructions, so the import plans to clear it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Import Clears Omitted"
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
}

data "lumos_current_user" "me" {}
