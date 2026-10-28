# requires: GROUP_ID
# The user is removed, leaving the group as the only admin.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Admins"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = var.group_id, type = "group" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
}

data "lumos_current_user" "me" {}
