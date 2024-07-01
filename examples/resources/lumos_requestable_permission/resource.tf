resource "lumos_requestable_permission" "aws_admin_permission" {
  app_id = local.aws_app_id
  label  = "Admin access to AWS"

  request_config = {
    appstore_visibility = "VISIBLE"

    request_approval_config = {
      request_approval_config_override = true
      manager_approval                 = "NONE"
      approvers = {
        users = [
          {
            id = local.aws_admin_user_id
          }
        ]
        groups = [
          {
            id = local.aws_admin_group_id
          }
        ]
      }
    }
    request_fulfillment_config = {

      time_based_access_override = true
      time_based_access          = ["2 hours", "4 hours"]
      provisioning_group = {
        integration_specific_id = "0g1g2g3g4g5g6g7g8g9"
        app_id                  = local.okta_app_id
      }
    }
  }
}
