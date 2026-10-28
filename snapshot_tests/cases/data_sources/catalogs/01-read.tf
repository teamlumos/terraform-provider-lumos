# Catalogues that are the same in every organization.
data "lumos_app_categories" "all" {}

data "lumos_integrations" "okta" {
  query = "okta"
}

data "lumos_integration_connection_schema" "okta" {
  app_class_id = "okta.com"
}

data "lumos_inline_webhooks" "all" {}

data "lumos_custom_attribute_definitions" "users" {
  attribute_type = "user"
}

output "categories_include_developers" {
  value = contains([for c in data.lumos_app_categories.all.app_categories : c.label], "Developers")
}

output "integrations_include_okta" {
  value = contains([for i in data.lumos_integrations.okta.integrations : i.app_class_id], "okta.com")
}

output "okta_auth_schema_is_json" {
  value = can(jsondecode(data.lumos_integration_connection_schema.okta.auth_schema))
}

output "webhook_types_are_known" {
  value = alltrue([for w in data.lumos_inline_webhooks.all.inline_webhooks : contains(["pre-approval", "provisioning", "deprovisioning", "request-validation", "siem"], w.type)])
}

output "user_filter_returns_only_users" {
  value = length(data.lumos_custom_attribute_definitions.users.custom_attribute_definitions) > 0 && alltrue([for d in data.lumos_custom_attribute_definitions.users.custom_attribute_definitions : d.attribute_type == "user"])
}
