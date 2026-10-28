# The second entitlement is pre-selected instead.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Default Entitlement"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "first" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

resource "lumos_entitlement" "second" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Write"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  default_entitlement                 = { id = lumos_entitlement.second.id }
}

data "lumos_current_user" "me" {}
