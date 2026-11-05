# app_class_id = "custom" added to the config, as when it's copied from
# `terraform show`. It matches what the app reads back.
resource "lumos_app" "this" {
  app_class_id         = "custom"
  label_override       = "bb-${var.run_id} Class ID Added"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
