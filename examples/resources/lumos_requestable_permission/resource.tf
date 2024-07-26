data "lumos_users" "aws_admin" {
  search_term = "admin"
}

data "lumos_apps" "aws_app" {
  name_search = "aws"
}

data "lumos_apps" "okta_app" {
  name_search = "okta"
}

resource "lumos_requestable_permission" "aws" {
  app_id = data.lumos_apps.aws_app.items[0].id
  label  = "aws-application"
  request_config = {
    appstore_visibility = "VISIBLE"
    request_approval_config = {
      request_approval_config_override = true
      approvers = {
        groups = []
        users = [
          {
            id = data.lumos_users.aws_admin.items[0].id
          },
        ]
      }
      approvers_stage_2 = {
        groups = []
        users  = []
      }
      manager_approval = "NONE"
    }
    request_fulfillment_config = {
      provisioning_group = {
        app_id                  = data.lumos_apps.okta_app.items[0].id
        integration_specific_id = "00gh39pjk109MXCdR697"
      }
      time_based_access = [
        "Unlimited",
      ]
      time_based_access_override = true
    }
  }
}
