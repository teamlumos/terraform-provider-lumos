# Same config as step 01: the plan creates the policy again, instead of
# failing to read it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Policy Drift"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_access_policy" "this" {
  name                   = "bb-${var.run_id} Policy Drift"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{ id = lumos_app.this.id }]
}
