resource "lumos_requestable_permission" "example_permission" {
  app_id = "b4bb3f9c-e287-4cda-878d-8638a85206f9"
  label  = "Permission Name"
  request_config = {
    appstore_visibility = "VISIBLE"
    request_approval_config = {
      manager_approval                 = "INITIAL_APPROVAL"
      request_approval_config_override = true
      approvers = {
        groups = []
        users  = [{ id = data.lumos_users.example_users.items[0].id }]
      }

      custom_approval_message          = "Custom Approval Message"
      custom_approval_message_override = true
    }
    request_fulfillment_config = {
      time_based_access_options  = ["4 hours", "12 hours"]
      time_based_access_override = true
      provisioning_group = {
        app_id                  = "dd68df37-55b1-ccfd-d82c-f53d1af1c676"
        integration_specific_id = "00gn10ldkhddQYp03647"
      }
    }
  }
}
