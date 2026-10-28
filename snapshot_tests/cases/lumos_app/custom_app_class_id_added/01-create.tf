# A custom app with its required fields. It reads back with
# app_class_id = "custom".
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Class ID Added"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
