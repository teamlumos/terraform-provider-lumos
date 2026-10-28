# A custom app and its App Store settings in one config. lumos_app shows the
# settings read-only in appstore_settings, so the plan after the apply shows
# whether the app's copy is planned again once the settings are written.
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
  request_instructions                = "Snapshot test instructions"
}

data "lumos_current_user" "me" {}
