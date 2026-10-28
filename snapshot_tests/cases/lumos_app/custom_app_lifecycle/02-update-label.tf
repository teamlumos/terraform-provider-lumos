# Only label_override changes, so only it and label should change in the plan.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deploy Dashboard v2"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
