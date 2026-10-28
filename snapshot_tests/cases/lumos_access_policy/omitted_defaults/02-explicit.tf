# The same values written out: false, false and an empty set. Nothing changes,
# so the plan is empty.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Policy Defaults"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Defaults"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false
  is_everyone_condition  = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{
    id             = lumos_app.this.id
    is_preapproved = false
    entitlements   = []
  }]
}
