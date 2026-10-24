variable "entitlement_id" {
  type = string
}

# Reads one entitlement by its ID, for example to show its label and the app
# it belongs to.
data "lumos_entitlement" "this" {
  id = var.entitlement_id
}

output "entitlement_label" {
  value = data.lumos_entitlement.this.label
}

output "entitlement_app_id" {
  value = data.lumos_entitlement.this.app_id
}
