# Same config as step 01: the plan creates the app again, instead of
# failing to read it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deleted In Lumos"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
