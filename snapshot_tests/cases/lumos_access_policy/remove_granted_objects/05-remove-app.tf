# The dashboard app leaves the config and the policy in one apply. As in step
# 02, Terraform deletes the app before it updates the policy.
resource "lumos_app" "wiki" {
  label_override       = "bb-${var.run_id} Granted Wiki"
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

  apps = [{ id = lumos_app.wiki.id }]
}
