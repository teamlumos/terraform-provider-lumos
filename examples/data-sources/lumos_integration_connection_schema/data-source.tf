# The auth and settings keys an integration takes when connected with
# lumos_app. Both are JSON Schema documents, as JSON strings.
data "lumos_integration_connection_schema" "okta" {
  app_class_id = "okta.com"
}

output "okta_credential_keys" {
  value = keys(jsondecode(data.lumos_integration_connection_schema.okta.auth_schema).properties)
}
