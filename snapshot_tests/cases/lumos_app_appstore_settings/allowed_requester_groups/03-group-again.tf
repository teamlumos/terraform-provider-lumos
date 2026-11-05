# requires: GROUP_ID
# The group again.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Requester Groups"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  allowed_requester_group_uuids       = [var.group_id]
}

data "lumos_current_user" "me" {}
