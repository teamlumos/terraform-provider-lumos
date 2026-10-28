# requires: GROUP_ID
# expect: error
# An admin id that isn't a user in the organization is refused.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Admins"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = "00000000-0000-4000-8000-000000000000", type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
}

data "lumos_current_user" "me" {}
