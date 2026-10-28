# A redirect_url of 2048 characters, the most the API allows.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Redirect Length"
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
  redirect_url                        = "https://example.com/${join("", [for i in range(1014) : "xx"])}"
}

data "lumos_current_user" "me" {}
