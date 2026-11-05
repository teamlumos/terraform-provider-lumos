# expect: error
# A custom app has no default description to fall back to, so leaving
# description_override out of its config is refused at apply.
resource "lumos_app" "this" {
  label_override    = "bb-${var.run_id} Required Fields"
  category_override = "Developers"
}
