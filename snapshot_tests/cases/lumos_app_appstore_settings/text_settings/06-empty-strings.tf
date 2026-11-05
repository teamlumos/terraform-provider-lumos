# expect: error
# Empty strings are refused at plan time; leaving them out clears them.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Text Settings"
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
  request_instructions                = ""
  redirect_url                        = ""
}

data "lumos_current_user" "me" {}
