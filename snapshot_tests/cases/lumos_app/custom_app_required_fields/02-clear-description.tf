# description_override is left out after create, which clears it: description
# shows the app's default.
resource "lumos_app" "this" {
  label_override    = "bb-${var.run_id} Required Fields"
  category_override = "Developers"
}
