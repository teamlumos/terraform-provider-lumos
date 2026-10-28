# Only the label changes, so only label should change in the plan.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Lifecycle"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id                           = lumos_app.this.id
  label                            = "bb-${var.run_id} Read Only"
  description_override             = "Read access to the dashboard."
  custom_approval_message_override = "Your read access is ready."
}
