# expect: error
# An empty description_override is refused at plan time.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Empty Label"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id               = lumos_app.this.id
  label                = "bb-${var.run_id} Read"
  description_override = ""
}
