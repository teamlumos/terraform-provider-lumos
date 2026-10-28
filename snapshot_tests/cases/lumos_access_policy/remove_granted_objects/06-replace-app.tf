# The wiki is swapped for a new app in one apply. Terraform creates the new
# app and deletes the wiki before it updates the policy.
resource "lumos_app" "docs" {
  label_override       = "bb-${var.run_id} Granted Docs"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
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

  apps = [{ id = lumos_app.docs.id }]
}
