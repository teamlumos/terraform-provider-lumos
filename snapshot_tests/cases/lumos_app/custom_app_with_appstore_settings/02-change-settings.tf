# Only the settings change. lumos_app should have nothing to do, now or in the
# plan after the apply.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} With Settings"
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
  request_instructions                = "Snapshot test instructions, edited"
}

data "lumos_current_user" "me" {}
