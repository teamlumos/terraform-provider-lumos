# Finds an app by name, for example to manage its App Store settings.
data "lumos_apps" "okta" {
  search = "Okta"
}

output "okta_app_id" {
  value = one([for app in data.lumos_apps.okta.apps : app.id if app.label == "Okta"])
}
