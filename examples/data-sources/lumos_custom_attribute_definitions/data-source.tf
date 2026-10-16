# The organization's custom attributes on apps.
data "lumos_custom_attribute_definitions" "app_attributes" {
  config_type = "DOMAIN_APP"
}

output "app_attribute_labels" {
  value = [for a in data.lumos_custom_attribute_definitions.app_attributes.custom_attribute_definitions : a.label]
}
