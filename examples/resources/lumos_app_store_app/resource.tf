resource "lumos_app_store_app" "my_appstoreapp" {
  app_id                      = "...my_app_id..."
  custom_request_instructions = "...my_custom_request_instructions..."
  provisioning = {
    access_removal_inline_webhook = {
      id = "...my_id..."
    }
    allow_multiple_permission_selection = true
    custom_provisioning_instructions    = "...my_custom_provisioning_instructions..."
    groups_provisioning                 = "GROUPS_AND_VISIBLE"
    manual_steps_needed                 = false
    provisioning_webhook = {
      id = "...my_id..."
    }
    time_based_access = [
      "90 days"
    ]
  }
  request_flow = {
    admins = {
      groups = [
        {
          app_id                  = "...my_app_id..."
          id                      = "...my_id..."
          integration_specific_id = "...my_integration_specific_id..."
        }
      ]
      users = [
        {
          id = "...my_id..."
        }
      ]
    }
    allowed_groups = {
      groups = [
        {
          app_id                  = "...my_app_id..."
          id                      = "...my_id..."
          integration_specific_id = "...my_integration_specific_id..."
        }
      ]
      type = "SPECIFIED_GROUPS"
    }
    approvers = {
      groups = [
        {
          app_id                  = "...my_app_id..."
          id                      = "...my_id..."
          integration_specific_id = "...my_integration_specific_id..."
        }
      ]
      users = [
        {
          id = "...my_id..."
        }
      ]
    }
    approvers_stage_2 = {
      groups = [
        {
          app_id                  = "...my_app_id..."
          id                      = "...my_id..."
          integration_specific_id = "...my_integration_specific_id..."
        }
      ]
      users = [
        {
          id = "...my_id..."
        }
      ]
    }
    custom_approval_message = "...my_custom_approval_message..."
    discoverability         = "LIMITED"
    request_validation_inline_webhook = {
      id = "...my_id..."
    }
    require_additional_approval = true
    require_manager_approval    = true
  }
}