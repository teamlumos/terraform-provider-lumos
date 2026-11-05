# A manual entitlement on a custom app, with a description and an approval
# message of its own.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Import"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id                           = lumos_app.this.id
  label                            = "bb-${var.run_id} Imported"
  description_override             = "Imported entitlement."
  custom_approval_message_override = "Your imported access is ready."
}
