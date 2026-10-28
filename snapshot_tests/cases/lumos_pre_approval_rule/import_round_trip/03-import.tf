# The rule is found through lumos_pre_approval_rules and imported with an import
# block and the config it was created with, so the import plans no changes.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Import"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "viewer" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Viewer"
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

data "lumos_pre_approval_rules" "this" {
  app_id = lumos_app.this.id
}

import {
  to = lumos_pre_approval_rule.this
  id = jsonencode({ app_id = lumos_app.this.id, id = one([for r in data.lumos_pre_approval_rules.this.pre_approval_rules : r.id if r.justification == "Snapshot test rule"]) })
}

resource "lumos_pre_approval_rule" "this" {
  app_id        = lumos_app.this.id
  justification = "Snapshot test rule"
  condition = jsonencode({
    and = [
      {
        attribute = { label = "Email", value = "identity.email" }
        equals    = { label = "bb-${var.run_id}@example.com", value = "bb-${var.run_id}@example.com" }
      },
      {
        attribute = { label = "Access length", value = "access_request.access_length" }
        "<="      = { label = "30 days", value = 2592000 }
      },
    ]
  })
  request_includes_entitlements = [lumos_entitlement.viewer.id]
}
