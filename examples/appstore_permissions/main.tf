terraform {
  required_providers {
    lumos-appstore = {
      source = "lumos.com/tf/lumos-appstore"
    }
  }
}

provider "lumos-appstore" {
  api_token = "<API TOKEN>"
}

// TODO Replace UUIDs with your local UUIDs For Testing
locals {
  terraform_test_app_id = "<FILL IN ID>" # Find UUID like so: https://www.loom.com/share/770d0a4900e7430f8ec8e552819eafdc?sid=61920d66-1a98-4b64-b8b1-844d422d6930
  group_1_id            = "<FILL IN ID>" # Find UUID like so: https://www.loom.com/share/b313034f7ad84d83b0195fe1b0f370f8?sid=ff1220e9-874c-4021-b915-d9be509ee1c2
  group_2_id            = "<FILL IN ID>"
  group_3_id            = "<FILL IN ID>"
  user_1_id             = "<FILL IN ID>" # Find UUID like so: https://www.loom.com/share/1413b71a5f8f4692a01f986a95b3e6fc?sid=5def1c2a-6cae-4c9b-b2d2-b944131ab5b8
  user_2_id             = "<FILL IN ID>"
  user_3_id             = "<FILL IN ID>"
}

// TODO Uncomment to create requestable permissions.
# resource "lumos-appstore_requestable_permission" "hello_terraform_one" {
#   app_id                            = local.terraform_test_app_id
#   label                             = "Hello Terraform One!"
#   visible_in_appstore               = false
#   allowed_groups                    = [local.group_1_id, local.group_2_id]
#   manager_approval_required         = false
#   require_additional_approval       = false
#   approver_groups_stage_1           = []
#   approver_users_stage_1            = [local.user_1_id]
#   approver_groups_stage_2           = []
#   approver_users_stage_2            = []
#   access_removal_inline_webhook     = null
#   request_validation_inline_webhook = null
#   provisioning_inline_webhook       = null
#   manual_steps_needed               = true
#   manual_instructions               = "It's very important you follow these instructions so that we can stay compliant. Thank you!"
#   time_based_access_options         = ["Unlimited"]
#   provisioning_group                 = local.group_3_id
# }

# resource "lumos-appstore_requestable_permission" "hello_terraform_two" {
#   app_id                            = local.terraform_test_app_id
#   label                             = "Hello Terraform Two!"
#   visible_in_appstore               = false
#   allowed_groups                    = null
#   manager_approval_required         = false
#   require_additional_approval       = true
#   approver_groups_stage_1           = []
#   approver_users_stage_1            = [local.user_2_id]
#   approver_groups_stage_2           = [local.group_1_id]
#   approver_users_stage_2            = []
#   access_removal_inline_webhook     = null
#   request_validation_inline_webhook = null
#   provisioning_inline_webhook       = null
#   manual_steps_needed               = true
#   manual_instructions               = "Please follow these critical steps!"
#   time_based_access_options         = ["Unlimited"]
#   provisioning_group                 = local.group_3_id
# }
