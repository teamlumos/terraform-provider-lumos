# Terraform stops managing the policy without deleting it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Policy Import"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

removed {
  from = lumos_access_policy.this

  lifecycle {
    destroy = false
  }
}
