# expect: error
# A policy needs access_condition or is_everyone_condition = true.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} No Condition"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} No Condition"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false

  apps = [{ id = lumos_app.this.id }]
}
