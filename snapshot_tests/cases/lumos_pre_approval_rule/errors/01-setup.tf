# Two apps and their entitlements for the refused rules that follow.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Errors"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "viewer" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Viewer"
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_app" "other" {
  label_override       = "bb-${var.run_id} Pre-Approval Other App"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "other" {
  app_id = lumos_app.other.id
  label  = "bb-${var.run_id} Other Read"
}
