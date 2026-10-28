# expect: error
# A draft can't be enabled.
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
  name                   = "bb-${var.run_id} Enabled Draft"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = true

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{ id = lumos_app.this.id }]
}
