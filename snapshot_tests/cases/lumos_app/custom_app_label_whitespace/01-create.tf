# A label with leading and trailing spaces. If the API trims it, the apply
# fails with an inconsistent result.
resource "lumos_app" "this" {
  label_override       = "  bb-${var.run_id} Spaced  "
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
