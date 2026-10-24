# Terraform manages who can request Okta. App owners keep editing the
# request instructions and redirect URL in the Lumos UI: Terraform doesn't
# plan to clear them and sends their current values on each update.
resource "lumos_app_appstore_settings" "okta" {
  app_id                              = lumos_app.okta.id
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  allowed_requester_group_uuids       = [var.engineering_group_id]

  lifecycle {
    ignore_changes = [request_instructions, redirect_url]
  }
}
