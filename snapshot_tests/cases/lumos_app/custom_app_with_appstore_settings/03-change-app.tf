# Only the app's label changes. Its appstore_settings shouldn't show as
# changing or unknown.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} With Settings v2"
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
