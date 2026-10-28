# A custom app.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deleted In Lumos"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
