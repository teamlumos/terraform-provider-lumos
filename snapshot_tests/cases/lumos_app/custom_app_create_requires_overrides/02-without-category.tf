# expect: error
# Creating a custom app without category_override is refused.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Requires Overrides"
  description_override = "Snapshot test custom app"
}
