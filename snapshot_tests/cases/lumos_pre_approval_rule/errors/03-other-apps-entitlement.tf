# expect: error
# An entitlement of another app in request_includes_entitlements.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Errors"
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

resource "lumos_app" "other" {
  label_override       = "bb-${var.run_id} Pre-Approval Other App"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "other" {
  app_id = lumos_app.other.id
  label  = "bb-${var.run_id} Other Read"
}

resource "lumos_pre_approval_rule" "this" {
  app_id        = lumos_app.this.id
  justification = "Snapshot test rule"
  condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    equals    = { label = "bb-${var.run_id}@example.com", value = "bb-${var.run_id}@example.com" }
  })
  request_includes_entitlements = [lumos_entitlement.other.id]
}
