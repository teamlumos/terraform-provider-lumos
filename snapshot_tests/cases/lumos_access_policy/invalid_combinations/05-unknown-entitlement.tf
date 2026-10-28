# expect: error
# An entitlement ID that doesn't exist.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Invalid Policies"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Unknown Entitlement"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{
    id           = lumos_app.this.id
    entitlements = [{ id = "00000000-0000-4000-8000-000000000000" }]
  }]
}
