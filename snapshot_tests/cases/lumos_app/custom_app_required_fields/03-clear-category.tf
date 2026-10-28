# expect: error
# The same for category_override: a custom app has no default category.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Required Fields"
  description_override = "Snapshot test custom app"
}
