# Only the settings the resource requires.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Webhook Wrong Type"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_inline_webhooks" "all" {}

locals {
  pre_approval_webhooks = [for w in data.lumos_inline_webhooks.all.inline_webhooks : w.id if w.type == "pre-approval"]
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  lifecycle {
    precondition {
      condition     = length(local.pre_approval_webhooks) > 0
      error_message = "This case needs a pre-approval inline webhook in the organization."
    }
  }
}

data "lumos_current_user" "me" {}
