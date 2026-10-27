# A custom app: one Lumos has no connector for, whose accounts you
# maintain yourself.
resource "lumos_app" "deploy_dashboard" {
  label_override       = "Deploy dashboard"
  category_override    = "Developers"
  description_override = "Internal tool for tracking deploys."
  website_url_override = "https://deploys.example.com"
}

# A connected integration. The keys auth and settings take depend on the
# integration; read them from its connection schema
# (GET /v1/integrations/{app_class_id}/connection-schema). label_override
# replaces the integration's name, as shown in Lumos.
variable "okta_api_token" {
  type      = string
  sensitive = true
}

resource "lumos_app" "okta" {
  app_class_id   = "okta.com"
  label_override = "Okta (Workforce)"
  auth = jsonencode({
    api_key = var.okta_api_token
  })
  settings = jsonencode({
    app_instance_identifier = "acme.okta.com"
  })
}
