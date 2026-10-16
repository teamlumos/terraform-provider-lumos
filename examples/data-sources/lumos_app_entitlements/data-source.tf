variable "app_id" {
  type = string
}

# One app's entitlements that are listed in the App Store.
data "lumos_app_entitlements" "requestable" {
  app_id              = var.app_id
  visible_in_appstore = true
}

output "requestable_labels" {
  value = [for e in data.lumos_app_entitlements.requestable.entitlements : e.label]
}
