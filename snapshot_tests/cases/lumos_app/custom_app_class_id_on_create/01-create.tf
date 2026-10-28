# A custom app created with app_class_id = "custom", the value a custom app
# reads back with.
resource "lumos_app" "this" {
  app_class_id         = "custom"
  label_override       = "bb-${var.run_id} Class ID On Create"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
