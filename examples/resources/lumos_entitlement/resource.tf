# A manual entitlement: one created in Lumos rather than synced from an
# integration. Its app can be a custom app or a connected integration.
resource "lumos_entitlement" "deploy_admin" {
  app_id = lumos_app.deploy_dashboard.id
  label  = "Admin"
}
