# The app is listed again, with the settings it had.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Required Settings"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "hidden"
  allow_permission_selection          = false
  allow_multiple_permission_selection = true
}

data "lumos_current_user" "me" {}
