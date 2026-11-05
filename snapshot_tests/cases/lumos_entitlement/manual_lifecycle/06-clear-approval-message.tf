# custom_approval_message_override is left out too, so the entitlement
# inherits the app's message.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Lifecycle"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read Only"
}
