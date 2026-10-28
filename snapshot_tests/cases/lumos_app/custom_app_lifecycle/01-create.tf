# A custom app with only the overrides a custom app requires on create.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deploy Dashboard"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
