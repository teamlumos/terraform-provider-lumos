# Terraform stops managing the rule without deleting it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Import"
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

removed {
  from = lumos_pre_approval_rule.this

  lifecycle {
    destroy = false
  }
}
