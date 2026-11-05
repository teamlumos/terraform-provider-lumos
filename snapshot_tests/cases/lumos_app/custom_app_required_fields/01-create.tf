# A custom app with the three overrides a custom app requires on create.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Required Fields"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
