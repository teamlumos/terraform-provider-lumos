# requires: LUMOS_TEST_TBA_LABEL
# Options and default set again.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Time Based Access"
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
  time_based_access = {
    options = [var.tba_label]
    default = var.tba_label
  }
}

data "lumos_current_user" "me" {}
