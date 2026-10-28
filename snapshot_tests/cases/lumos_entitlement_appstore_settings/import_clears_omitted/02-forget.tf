# requires: GROUP_ID
# Terraform stops managing the settings without changing them.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Import Clears"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = true
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

removed {
  from = lumos_entitlement_appstore_settings.this

  lifecycle {
    destroy = false
  }
}

data "lumos_current_user" "me" {}
