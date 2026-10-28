# The org's first request-validation webhook validates requests for the
# entitlement.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Webhook"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = true
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

data "lumos_inline_webhooks" "all" {}

locals {
  validation_webhooks   = [for w in data.lumos_inline_webhooks.all.inline_webhooks : w.id if w.type == "request-validation"]
  pre_approval_webhooks = [for w in data.lumos_inline_webhooks.all.inline_webhooks : w.id if w.type == "pre-approval"]
}

resource "lumos_entitlement_appstore_settings" "this" {
  entitlement_id                  = lumos_entitlement.this.id
  visible_in_appstore             = true
  request_validation_webhook_uuid = try(local.validation_webhooks[0], null)
  lifecycle {
    precondition {
      condition     = length(local.validation_webhooks) > 0
      error_message = "This case needs a request-validation inline webhook in the organization."
    }
  }

  depends_on = [lumos_app_appstore_settings.this]
}

data "lumos_current_user" "me" {}
