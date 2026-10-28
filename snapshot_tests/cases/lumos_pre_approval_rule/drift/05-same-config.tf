# Same config as step 01: the plan creates the rule again, instead of failing
# to read it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Drift"
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

resource "lumos_pre_approval_rule" "this" {
  app_id        = lumos_app.this.id
  justification = "Snapshot test rule"
  condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    equals    = { label = "bb-${var.run_id}@example.com", value = "bb-${var.run_id}@example.com" }
  })
  request_includes_entitlements = [lumos_entitlement.viewer.id]
}
