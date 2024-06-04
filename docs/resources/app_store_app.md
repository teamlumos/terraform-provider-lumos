---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lumos_app_store_app Resource - terraform-provider-lumos"
subcategory: ""
description: |-
  AppStoreApp Resource
---

# lumos_app_store_app (Resource)

AppStoreApp Resource

## Example Usage

```terraform
resource "lumos_app_store_app" "my_appstoreapp" {
  app_id                      = "...my_app_id..."
  custom_request_instructions = "...my_custom_request_instructions..."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_id` (String) The ID of the app to add to the app store. Requires replacement if changed.

### Optional

- `custom_request_instructions` (String) AppStore App instructions. Requires replacement if changed.
- `provisioning` (Attributes) Provisioning flow configuration to request access to app. Requires replacement if changed. (see [below for nested schema](#nestedatt--provisioning))
- `request_flow` (Attributes) Request flow configuration to request access to app. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow))

### Read-Only

- `allow_multiple_permission_selection` (Boolean) Whether the app is configured to allow multiple permissions to be requested at a time.
- `app_class_id` (String) The ID of the service associated with this app.
- `id` (String) The ID of this app.
- `instance_id` (String) The ID of the instance associated with this app.
- `sources` (List of String) The sources of this app.
- `status` (String) An enumeration. must be one of ["DISCOVERED", "NEEDS_REVIEW", "APPROVED", "BLOCKLISTED", "DEPRECATED"]
- `user_friendly_label` (String) The user-friendly label of this app.

<a id="nestedatt--provisioning"></a>
### Nested Schema for `provisioning`

Optional:

- `access_removal_inline_webhook` (Attributes) An inactivity workflow can be optionally associated with this app. Requires replacement if changed. (see [below for nested schema](#nestedatt--provisioning--access_removal_inline_webhook))
- `custom_provisioning_instructions` (String) Only Available if manual steps is active. During the provisioning step, send a custom message to app admins explaining how to provision a user to the app. Markdown for links and text formatting is supported. Requires replacement if changed.
- `groups_provisioning` (String) An enumeration. Requires replacement if changed. ; must be one of ["DIRECT_TO_USER", "GROUPS_AND_HIDDEN", "GROUPS_AND_VISIBLE"]
- `manual_steps_needed` (Boolean) If enabled, Lumos will reach out to the App Admin after initial access is granted to perform additional manual steps. Note that if this option is enabled, this action must be confirmed by the App Admin in order to resolve the request. Requires replacement if changed.
- `provisioning_webhook` (Attributes) The provisioning webhook optionally associated with this app. Requires replacement if changed. (see [below for nested schema](#nestedatt--provisioning--provisioning_webhook))
- `time_based_access` (List of String) If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access. Requires replacement if changed.

<a id="nestedatt--provisioning--access_removal_inline_webhook"></a>
### Nested Schema for `provisioning.access_removal_inline_webhook`

Optional:

- `id` (String) The ID of this inline webhook. Requires replacement if changed. ; Not Null

Read-Only:

- `description` (String) The description of this inline webhook.
- `hook_type` (String) An enumeration. must be one of ["PRE_APPROVAL", "PROVISION", "DEPROVISION", "REQUEST_VALIDATION", "SIEM"]
- `name` (String) The name of this inline webhook.


<a id="nestedatt--provisioning--provisioning_webhook"></a>
### Nested Schema for `provisioning.provisioning_webhook`

Optional:

- `id` (String) The ID of this inline webhook. Requires replacement if changed. ; Not Null

Read-Only:

- `description` (String) The description of this inline webhook.
- `hook_type` (String) An enumeration. must be one of ["PRE_APPROVAL", "PROVISION", "DEPROVISION", "REQUEST_VALIDATION", "SIEM"]
- `name` (String) The name of this inline webhook.



<a id="nestedatt--request_flow"></a>
### Nested Schema for `request_flow`

Optional:

- `admins` (Attributes) AppStore App admins assigned. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--admins))
- `allowed_groups` (Attributes) The allowed groups associated with this config. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--allowed_groups))
- `approvers` (Attributes) AppStore App approvers assigned. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--approvers))
- `approvers_stage_2` (Attributes) AppStore App stage 2 approvers assigned. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--approvers_stage_2))
- `custom_approval_message` (String) During the approval step, send a custom message to requesters. Markdown for links and text formatting is supported. Requires replacement if changed.
- `discoverability` (String) An enumeration. Requires replacement if changed. ; must be one of ["FULL", "LIMITED", "NONE"]
- `request_validation_inline_webhook` (Attributes) A request validation webhook can be optionally associated with this app. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--request_validation_inline_webhook))
- `require_additional_approval` (Boolean) Only turn on when working with sensitive permissions to ensure a smooth employee experience. Requires replacement if changed.
- `require_manager_approval` (Boolean) When a user makes an access request, require that their manager approves the request before moving on to additional approvals. Requires replacement if changed.

<a id="nestedatt--request_flow--admins"></a>
### Nested Schema for `request_flow.admins`

Optional:

- `groups` (Attributes List) Groups assigned as app admins. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--admins--groups))
- `users` (Attributes List) Users assigned as app admins. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--admins--users))

<a id="nestedatt--request_flow--admins--groups"></a>
### Nested Schema for `request_flow.admins.groups`

Optional:

- `app_id` (String) The ID of the app that owns this group. Requires replacement if changed.
- `id` (String) The ID of this group. Requires replacement if changed.
- `integration_specific_id` (String) The ID of this group, specific to the integration. Requires replacement if changed.

Read-Only:

- `description` (String) The description of this group.
- `group_lifecycle` (String) The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]
- `name` (String) The name of this group.
- `source_app_id` (String) The ID of the app that owns this group.


<a id="nestedatt--request_flow--admins--users"></a>
### Nested Schema for `request_flow.admins.users`

Optional:

- `id` (String) The ID of this user. Requires replacement if changed. ; Not Null

Read-Only:

- `email` (String) The email of this user.
- `family_name` (String) The family name of this user.
- `given_name` (String) The given name of this user.
- `status` (String) An enumeration. must be one of ["STAGED", "ACTIVE", "SUSPENDED", "INACTIVE"]



<a id="nestedatt--request_flow--allowed_groups"></a>
### Nested Schema for `request_flow.allowed_groups`

Optional:

- `groups` (Attributes List) The groups associated with this config. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--allowed_groups--groups))
- `type` (String) The type of this allowed groups config, can be all groups or specific. Requires replacement if changed. ; must be one of ["ALL_GROUPS", "SPECIFIED_GROUPS"]; Default: "ALL_GROUPS"

<a id="nestedatt--request_flow--allowed_groups--groups"></a>
### Nested Schema for `request_flow.allowed_groups.groups`

Optional:

- `app_id` (String) The ID of the app that owns this group. Requires replacement if changed.
- `id` (String) The ID of this group. Requires replacement if changed.
- `integration_specific_id` (String) The ID of this group, specific to the integration. Requires replacement if changed.

Read-Only:

- `description` (String) The description of this group.
- `group_lifecycle` (String) The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]
- `name` (String) The name of this group.
- `source_app_id` (String) The ID of the app that owns this group.



<a id="nestedatt--request_flow--approvers"></a>
### Nested Schema for `request_flow.approvers`

Optional:

- `groups` (Attributes List) Groups assigned as support request approvers. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--approvers--groups))
- `users` (Attributes List) Users assigned as support request approvers. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--approvers--users))

<a id="nestedatt--request_flow--approvers--groups"></a>
### Nested Schema for `request_flow.approvers.groups`

Optional:

- `app_id` (String) The ID of the app that owns this group. Requires replacement if changed.
- `id` (String) The ID of this group. Requires replacement if changed.
- `integration_specific_id` (String) The ID of this group, specific to the integration. Requires replacement if changed.

Read-Only:

- `description` (String) The description of this group.
- `group_lifecycle` (String) The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]
- `name` (String) The name of this group.
- `source_app_id` (String) The ID of the app that owns this group.


<a id="nestedatt--request_flow--approvers--users"></a>
### Nested Schema for `request_flow.approvers.users`

Optional:

- `id` (String) The ID of this user. Requires replacement if changed. ; Not Null

Read-Only:

- `email` (String) The email of this user.
- `family_name` (String) The family name of this user.
- `given_name` (String) The given name of this user.
- `status` (String) An enumeration. must be one of ["STAGED", "ACTIVE", "SUSPENDED", "INACTIVE"]



<a id="nestedatt--request_flow--approvers_stage_2"></a>
### Nested Schema for `request_flow.approvers_stage_2`

Optional:

- `groups` (Attributes List) Groups assigned as support request approvers. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--approvers_stage_2--groups))
- `users` (Attributes List) Users assigned as support request approvers. Requires replacement if changed. (see [below for nested schema](#nestedatt--request_flow--approvers_stage_2--users))

<a id="nestedatt--request_flow--approvers_stage_2--groups"></a>
### Nested Schema for `request_flow.approvers_stage_2.groups`

Optional:

- `app_id` (String) The ID of the app that owns this group. Requires replacement if changed.
- `id` (String) The ID of this group. Requires replacement if changed.
- `integration_specific_id` (String) The ID of this group, specific to the integration. Requires replacement if changed.

Read-Only:

- `description` (String) The description of this group.
- `group_lifecycle` (String) The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]
- `name` (String) The name of this group.
- `source_app_id` (String) The ID of the app that owns this group.


<a id="nestedatt--request_flow--approvers_stage_2--users"></a>
### Nested Schema for `request_flow.approvers_stage_2.users`

Optional:

- `id` (String) The ID of this user. Requires replacement if changed. ; Not Null

Read-Only:

- `email` (String) The email of this user.
- `family_name` (String) The family name of this user.
- `given_name` (String) The given name of this user.
- `status` (String) An enumeration. must be one of ["STAGED", "ACTIVE", "SUSPENDED", "INACTIVE"]



<a id="nestedatt--request_flow--request_validation_inline_webhook"></a>
### Nested Schema for `request_flow.request_validation_inline_webhook`

Optional:

- `id` (String) The ID of this inline webhook. Requires replacement if changed. ; Not Null

Read-Only:

- `description` (String) The description of this inline webhook.
- `hook_type` (String) An enumeration. must be one of ["PRE_APPROVAL", "PROVISION", "DEPROVISION", "REQUEST_VALIDATION", "SIEM"]
- `name` (String) The name of this inline webhook.

## Import

Import is supported using the following syntax:

```shell
terraform import lumos_app_store_app.my_lumos_app_store_app ""
```