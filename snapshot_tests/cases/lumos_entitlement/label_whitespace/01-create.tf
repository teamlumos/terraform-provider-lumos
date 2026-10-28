# A label with leading and trailing spaces. If the API trims it, the apply
# fails with an inconsistent result.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Whitespace"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "  bb-${var.run_id} Spaced  "
}
