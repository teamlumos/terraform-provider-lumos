# The condition is removed and the draft applies to everyone again.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Everyone Policy"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Everyone"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false
  is_everyone_condition  = true

  apps = [{ id = lumos_app.this.id }]
}
