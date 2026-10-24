variable "app_id" {
  type = string
}

# Reads one app by its ID, for example to show its name and whether it's
# listed in the App Store.
data "lumos_app" "this" {
  id = var.app_id
}

output "app_label" {
  value = data.lumos_app.this.label
}

output "app_in_app_store" {
  value = data.lumos_app.this.appstore_settings.in_app_store
}
