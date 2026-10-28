# A custom app with its required fields.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Version"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
