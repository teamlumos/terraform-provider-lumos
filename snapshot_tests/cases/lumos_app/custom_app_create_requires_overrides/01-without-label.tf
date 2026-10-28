# expect: error
# Creating a custom app without label_override is refused.
resource "lumos_app" "this" {
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
