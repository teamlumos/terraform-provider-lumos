# Terraform stops managing the settings without changing them.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Import Clears Omitted"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

removed {
  from = lumos_app_appstore_settings.this

  lifecycle {
    destroy = false
  }
}
