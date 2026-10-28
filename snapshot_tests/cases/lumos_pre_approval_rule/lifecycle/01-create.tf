# A rule pre-approving the Viewer entitlement for one run-specific email
# address, which no one has, so it pre-approves nothing.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Lifecycle"
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
