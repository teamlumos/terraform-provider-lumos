# description_override is left out, which clears it. A manual entitlement has
# no integration description to fall back to, so description is cleared too.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Lifecycle"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id                           = lumos_app.this.id
  label                            = "bb-${var.run_id} Read Only"
  custom_approval_message_override = "Your read access is ready. It can take a few minutes to appear."
}
