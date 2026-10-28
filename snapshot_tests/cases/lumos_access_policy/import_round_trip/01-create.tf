# A published policy granting an app and an entitlement.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Policy Import"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Imported"
  business_justification = "Snapshot test policy"
  status                 = "published"
  is_enabled             = true

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{
    id             = lumos_app.this.id
    is_preapproved = true
    entitlements   = [{ id = lumos_entitlement.this.id }]
  }]
}
