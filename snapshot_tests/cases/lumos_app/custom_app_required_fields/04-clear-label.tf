# expect: error
# The same for label_override: a custom app has no default name.
resource "lumos_app" "this" {
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
