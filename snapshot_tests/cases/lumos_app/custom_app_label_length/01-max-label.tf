# A label of 255 characters, the most the API allows.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id}-${join("", [for i in range(245) : "x"])}"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
