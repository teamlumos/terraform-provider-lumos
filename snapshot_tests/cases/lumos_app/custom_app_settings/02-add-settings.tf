# expect: error
# settings only applies to integrations, so the API refuses it on a custom app.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Settings"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
  settings             = jsonencode({})
}
