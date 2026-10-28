# The app leaves the App Store; Terraform keeps managing its settings.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Required Settings"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = false
  catalog_visibility                  = "hidden"
  allow_permission_selection          = false
  allow_multiple_permission_selection = true
}
