# The Admin entitlement leaves the config and the policy in one apply, and the
# policy is updated before the entitlement is deleted.
resource "lumos_app" "wiki" {
  label_override       = "bb-${var.run_id} CBD Wiki"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "dashboard" {
  label_override       = "bb-${var.run_id} CBD Dashboard"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"

  lifecycle {
    create_before_destroy = true
  }
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} CBD Granted Objects"
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
