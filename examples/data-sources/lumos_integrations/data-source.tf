# Finds an integration's app_class_id, which lumos_app uses to connect it.
data "lumos_integrations" "okta" {
  query = "Okta"
}

output "okta_app_class_ids" {
  value = [for i in data.lumos_integrations.okta.integrations : i.app_class_id]
}
