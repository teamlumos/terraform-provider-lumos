# A rule that pre-approves any of the app's entitlements when the org's first
# pre-approval webhook approves the request.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Webhook"
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

data "lumos_inline_webhooks" "all" {}

locals {
  pre_approval_webhooks = [for w in data.lumos_inline_webhooks.all.inline_webhooks : w.id if w.type == "pre-approval"]
}

resource "lumos_pre_approval_rule" "this" {
  app_id        = lumos_app.this.id
  justification = "Snapshot test webhook rule"
  condition = jsonencode({
    attribute = { label = "Access check", value = "webhook.responses[\"${try(local.pre_approval_webhooks[0], "")}\"]" }
    equals    = { label = "Pre-approved", value = "PREAPPROVED" }
  })
  request_has_entitlements = true

  lifecycle {
    precondition {
      condition     = length(local.pre_approval_webhooks) > 0
      error_message = "This case needs a pre-approval inline webhook in the organization."
    }
  }
}
