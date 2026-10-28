# The org's first request-validation webhook validates requests.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Validation Webhook"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_inline_webhooks" "all" {}

locals {
  validation_webhooks = [for w in data.lumos_inline_webhooks.all.inline_webhooks : w.id if w.type == "request-validation"]
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  request_validation_webhook_uuid     = try(local.validation_webhooks[0], null)
  lifecycle {
    precondition {
      condition     = length(local.validation_webhooks) > 0
      error_message = "This case needs a request-validation inline webhook in the organization."
    }
  }
}

data "lumos_current_user" "me" {}
