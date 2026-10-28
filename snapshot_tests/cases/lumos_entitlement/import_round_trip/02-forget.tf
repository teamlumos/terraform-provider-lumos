# Terraform stops managing the entitlement without deleting it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Import"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

removed {
  from = lumos_entitlement.this

  lifecycle {
    destroy = false
  }
}
