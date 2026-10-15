variable "engineering_group_id" {
  type = string
}

# Lists the entitlement in the App Store for one group, with its own
# time-based access options. An override omitted from the config inherits
# the app's setting.
resource "lumos_entitlement_appstore_settings" "deploy_admin" {
  entitlement_id                         = lumos_entitlement.deploy_admin.id
  visible_in_appstore                    = true
  allowed_requester_group_uuids_override = [var.engineering_group_id]

  time_based_access_override = {
    options = ["1 day"]
  }
}
