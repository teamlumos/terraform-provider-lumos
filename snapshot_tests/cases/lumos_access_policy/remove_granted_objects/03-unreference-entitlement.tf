# The two-apply way round step 02: first the policy stops granting the Admin
# entitlement, which stays in the config.
resource "lumos_app" "wiki" {
  label_override       = "bb-${var.run_id} Granted Wiki"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "dashboard" {
  label_override       = "bb-${var.run_id} Granted Dashboard"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.dashboard.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Granted Objects"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [
    { id = lumos_app.wiki.id },
    { id = lumos_app.dashboard.id },
  ]
}
