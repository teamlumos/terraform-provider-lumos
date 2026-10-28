# Two custom apps with the same label, in one apply.
resource "lumos_app" "first" {
  label_override       = "bb-${var.run_id} Same Label"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "second" {
  label_override       = "bb-${var.run_id} Same Label"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
