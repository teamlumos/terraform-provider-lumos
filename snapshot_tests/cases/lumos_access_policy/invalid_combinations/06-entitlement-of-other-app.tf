# expect: error
# An entitlement listed under an app it doesn't belong to.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Invalid Policies"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

resource "lumos_app" "other" {
  label_override       = "bb-${var.run_id} Invalid Policies Other"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Wrong App"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{
    id           = lumos_app.other.id
    entitlements = [{ id = lumos_entitlement.this.id }]
  }]
}
