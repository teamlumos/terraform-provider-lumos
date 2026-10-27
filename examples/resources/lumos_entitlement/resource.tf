# A manual entitlement: one created in Lumos rather than synced from an
# integration. Its app can be a custom app or a connected integration.
resource "lumos_entitlement" "deploy_admin" {
  app_id = lumos_app.deploy_dashboard.id
  label  = "Admin"

  description_override             = "Full admin access to Deploy Dashboard. Ask in #deploys before requesting."
  custom_approval_message_override = "Your Deploy Dashboard admin access is approved. Ask in #deploys if you need help."
}
