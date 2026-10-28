# Two apps; the entitlement is on the first.
resource "lumos_app" "first" {
  label_override       = "bb-${var.run_id} Move Entitlement, first"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "second" {
  label_override       = "bb-${var.run_id} Move Entitlement, second"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.first.id
  label  = "bb-${var.run_id} Movable"
}
