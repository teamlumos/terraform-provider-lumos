# The policy is disabled and stays published.
resource "lumos_app" "wiki" {
  label_override       = "bb-${var.run_id} Policy Wiki"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "dashboard" {
  label_override       = "bb-${var.run_id} Policy Dashboard"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.dashboard.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Policy Renamed"
  business_justification = "Snapshot test policy, edited"
  status                 = "published"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [
    {
      id             = lumos_app.wiki.id
      is_preapproved = true
    },
    { id = lumos_app.dashboard.id },
  ]
}
