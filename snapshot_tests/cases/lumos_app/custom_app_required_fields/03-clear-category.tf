# category_override is cleared too: category shows the app's default.
resource "lumos_app" "this" {
  label_override = "bb-${var.run_id} Required Fields"
}
