# expect: error
# A label of 256 characters is refused, and the app keeps its label.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id}-${join("", [for i in range(246) : "x"])}"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
