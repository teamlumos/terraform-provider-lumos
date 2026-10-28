# request_does_not_include_entitlements replaces it: the rule covers every
# entitlement but Admin, and a request for the app itself.
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
  justification = "Snapshot test rule, edited"
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
  request_does_not_include_entitlements = [lumos_entitlement.admin.id]
}
