# expect: error
# The pre-selected entitlement must belong to this app.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Default Other App"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "other" {
  label_override       = "bb-${var.run_id} Default Other App, other"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "other" {
  app_id = lumos_app.other.id
  label  = "bb-${var.run_id} Other"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  default_entitlement                 = { id = lumos_entitlement.other.id }
}

data "lumos_current_user" "me" {}
