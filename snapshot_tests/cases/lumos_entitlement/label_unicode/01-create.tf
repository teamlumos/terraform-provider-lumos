# A label with accented letters and an emoji.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Unicode"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Accès 🔑"
}
