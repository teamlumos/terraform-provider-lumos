---
page_title: "lumos_appstore-requestable_permission Resource - Lumos AppStore Early Access Provider"
description: |-
  AppStore Permissions for an Application.
---

# lumos_appstore-requestable_permission (Resource)

## Example Usage

```terraform
resource "lumos_appstore-requestable_permission" "example" {
  app_id                            = "80368729-3603-bc74-a8ab-90b8abaebc29"
  label                             = "Example Permission"
  visible_in_appstore               = true
  allowed_groups                    = null
  manager_approval_required         = false
  require_additional_approval       = false
  approver_groups_stage_1           = []
  approver_users_stage_1            = ["75d511aa-5e13-eae1-3406-add39b321f62]
  approver_groups_stage_2           = []
  approver_users_stage_2            = []
  access_removal_inline_webhook     = null
  request_validation_inline_webhook = "9e078df2-3dff-668b-5647-9ac8eac6c7ad"
  provisioning_inline_webhook       = null
  manual_steps_needed               = true
  manual_instructions               = "Instructions for manual steps are here."
  time_based_access_options         = ["Unlimited"]
  provisioning_group                 = "b9e77977-4846-3d2a-087a-5a089b176622"
}
```

## Schema

- `app_id` (String UUID) The ID of the lumos app associated with this requestable permission.
- `label` (String) The label of this requestable permission.
- `visible_in_appstore` (Boolean) Defines if the permission will be visible in the AppStore to request access or not.
- `manager_approval_required` (Boolean) If enabled, when a user makes an access request fo this permission, it will require that their manager approves the request before moving on to additional approvals.
- `require_additional_approval` (Boolean) Only turn on when working with sensitive permissions to ensure a smooth employee experience.
- `time_based_access_options` (String[]) Access length options available for the users to request an app for a selected duration. After expiry, Lumos will automatically remove user's access. If empty default will be Unlimited.
  - Available options: [”2 hours”, “4 hours”, “12 hours”,”7 days”, “14 days”, “30 days”, “90 days”, “Unlimited”]
  - This options could be restricted based on domain configured as available options.
  - Permissions could have TBA options defined different than Unlimited only when they have access removal webhook attached and/or provisioning group.
- `manual_steps_needed` (Boolean) If enabled, Lumos will reach out to the App Admin after initial access is granted to perform additional manual steps. Note that if this option is enabled, this action must be confirmed by the App Admin in order to resolve the request.

### Optional

- `allowed_groups` (String UUID[]) A list of allowed groups to request access to this permission.
- `approver_groups_stage_1` (String UUID[]) AppStore App group IDs assigned as approvers.
- `approver_users_stage_1` (String UUID[]) AppStore App user IDs assigned as approvers.
- `approver_groups_stage_2` (String UUID[]) AppStore App group IDs assigned as approvers stage2.
- `approver_users_stage_2` (String UUID[]) AppStore App user IDs assigned as approvers stage2.
- `request_validation_inline_webhook` (String UUID) The UUID of request validation webhook can be associated with this permission.
- `provisioning_inline_webhook` (String UUID) The UUID of provisioning webhook associated with this permission.
- `access_removal_inline_webhook` (String UUID) The ID of inactivity workflow can be associated with this permission.
- `manual_instructions` (String) Only Available if manual steps is active. During the provisioning step, send a custom message to app admins explaining how to provision a user to the permission. Markdown for links and text formatting is supported.
- `provisioning_group` (String UUID) The group UUID associated with this config.

### Read-Only

- `id` (String) Requestable permission identifier
- `permission_type` (String) Defines if its a "SYNCED" permission from idp with a real association with the app or a "NATIVE" [manual] permission. This field will become required the next iteration, for now all permissions will be "NATIVE".
