# Two apps; the settings list the first one.
resource "lumos_app" "first" {
  label_override       = "bb-${var.run_id} Move Settings, first"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "second" {
  label_override       = "bb-${var.run_id} Move Settings, second"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.first.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  request_instructions                = "Snapshot test instructions"
}

data "lumos_current_user" "me" {}
