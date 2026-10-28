# expect: error
# Creating a custom app without description_override is refused.
resource "lumos_app" "this" {
  label_override    = "bb-${var.run_id} Requires Overrides"
  category_override = "Developers"
}
