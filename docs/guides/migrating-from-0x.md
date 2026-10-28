---
page_title: "Migrating from 0.x to 1.x"
subcategory: ""
description: |-
  Move a Terraform configuration from the 0.x provider to the 1.x provider.
---

# Migrating from 0.x to 1.x

Version 1.x of the provider is built on the Lumos V1 API. Some resources were renamed or split, App Store approvals and provisioning moved into workflows, and permissions are managed directly as the entitlements behind them, so they're addressed by the entitlement's ID instead of the requestable permission's ID. Terraform can't carry 0.x state over to the new resource types, so migrating means removing the 0.x resources from state without touching Lumos, then importing the same objects as 1.x resources.

The migration itself doesn't create, change or delete anything in Lumos. The first `terraform apply` on 1.x should show no changes other than the ones listed under [What to expect on the first apply](#what-to-expect-on-the-first-apply).

## Before you start

* **The V1 API must be enabled for your organization.**
* **Terraform 1.7 or later** for `removed` blocks and `import` blocks with `for_each`. Terraform 1.5 works if you use `terraform state rm` and write one `import` block per resource.
* **Back up your state:** `terraform state pull > pre-1x.tfstate`.
* **Start from a clean plan on 0.x.** Run `terraform plan` on your current 0.x version and resolve any changes first, so the migration doesn't mix with unrelated drift.

## What changed

### Resources

| 0.x | 1.x | Same ID? |
|---|---|---|
| `lumos_app` | `lumos_app` | Yes |
| `lumos_app_store_app` | `lumos_app_appstore_settings`, plus approval and provisioning workflows (see below) | Yes: the app's ID |
| `lumos_requestable_permission` | `lumos_entitlement` + `lumos_entitlement_appstore_settings`, plus workflows for approval and provisioning overrides | **No** (see [Finding entitlement IDs](#finding-entitlement-ids)) |
| `lumos_access_policy` | `lumos_access_policy` | Yes |
| `lumos_pre_approval_rule` | `lumos_pre_approval_rule` | Yes |
| Approvers, manager approval, stage 2 approvers, provisioning webhooks, manual provisioning steps (fields on `lumos_app_store_app` and `lumos_requestable_permission`) | `lumos_workflow` + a binding resource: `lumos_app_approval_workflow`, `lumos_app_provisioning_workflow`, `lumos_app_deprovisioning_workflow`, `lumos_entitlement_approval_workflow`, `lumos_entitlement_provisioning_workflow`, `lumos_entitlement_deprovisioning_workflow` | New resources. Import the workflows Lumos already has |

### Data sources

| 0.x | 1.x |
|---|---|
| `lumos_app` (`id`) | `lumos_app` (`id`) |
| `lumos_apps` | `lumos_apps` (see [`lumos_apps` returns every app](#lumos_apps-returns-every-app)) |
| `lumos_app_store_app`, `lumos_app_store_app_settings` | `appstore_settings` on `lumos_app` or `lumos_apps` |
| `lumos_requestable_permission` (`id`) | `lumos_entitlement` (`id`, the entitlement's ID) |
| `lumos_requestable_permissions` | `lumos_app_entitlements` (one app) or `lumos_entitlements` (every app; `search` required) |
| `lumos_user` (`user_id`) | `lumos_user` (`id`), plus `lumos_current_user` |
| `lumos_users` | `lumos_users` (`search` required) |
| `lumos_group` (`group_id`) | `lumos_group` (`id`) |
| `lumos_groups` | `lumos_groups` (`name`, `integration_specific_id` or `app_id` required) |
| `lumos_access_policy`, `lumos_access_policies` | Same names |
| `lumos_pre_approval_rule` (`id`) | `lumos_pre_approval_rule` (`app_id` and `id`), plus `lumos_pre_approval_rules` |

List data sources page through every result for you. `page`, `size`, `total` and `pages` are gone, and results sit under a named attribute instead of `items`, for example `data.lumos_apps.x.apps` and `data.lumos_users.x.users`. `name_search` and `search_term` are now `search`.

New lookups: `lumos_integrations`, `lumos_integration_connection_schema`, `lumos_inline_webhooks`, `lumos_app_categories`, `lumos_custom_attribute_definitions`, `lumos_accounts`, `lumos_user_roles`, `lumos_workflow_steps`.

### Values are lowercase

Enum values use the V1 API's lowercase, hyphenated spelling: `status = "published"` instead of `"PUBLISHED"`, `connection_source = "api"` instead of `"API"`, `catalog_visibility = "search-only"`.

### Leaving a setting out clears it

In 0.x, a setting left out of the configuration kept its value in Lumos. In 1.x, the configuration is authoritative: a setting left out, or set to `null`, is cleared, including values set in the Lumos UI, and the plan shows it as `-> null` or `-> []`. This applies to `lumos_app_appstore_settings`, `lumos_entitlement_appstore_settings`, and the `*_override` fields and `post_approval_message` on `lumos_app` and `lumos_entitlement`.

The steps below generate the 1.x configuration from what's in Lumos, rather than translating the 0.x configuration by hand, for this reason. To keep editing some settings in the Lumos UI, see [Leaving settings to the Lumos UI](#leaving-settings-to-the-lumos-ui).

## Migration steps

### 1. List what Terraform manages

```shell
terraform state list | grep '^lumos_'
```

Note each resource's ID (`terraform state show <address>`). For `lumos_app_store_app`, the ID to import is its `app_id`.

### 2. Stop managing the 0.x resources, without changing Lumos

Replace each 0.x `resource` block with a `removed` block, then apply while still on 0.x:

```terraform
removed {
  from = lumos_requestable_permission.aws_admin
  lifecycle {
    destroy = false
  }
}
```

```shell
terraform apply
```

The plan lists each resource as no longer managed and destroys nothing. Check that before applying. On Terraform older than 1.7, run `terraform state rm <address>` for each resource instead.

Then delete the `removed` blocks and any 0.x `lumos_*` data source blocks.

### 3. Switch to 1.x

```terraform
terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "~> 1.0"
    }
  }
}
```

```shell
terraform init -upgrade
```

The provider configuration (`http_bearer` or `LUMOS_ACCESS_TOKEN`, and `server_url`) is unchanged.

### 4. Write `import` blocks

Apps, access policies and pre-approval rules keep their 0.x IDs:

```terraform
import {
  to = lumos_app.okta
  id = "<app id>"
}

# Was lumos_app_store_app.okta. The import ID is the app's ID.
import {
  to = lumos_app_appstore_settings.okta
  id = "<app id>"
}

import {
  to = lumos_access_policy.engineering
  id = "<access policy id>"
}

import {
  to = lumos_pre_approval_rule.platform_viewers
  id = "<pre-approval rule id>"
}
```

Each `lumos_requestable_permission` is imported as the `lumos_entitlement` behind it, using the entitlement's ID (see [Finding entitlement IDs](#finding-entitlement-ids)), plus a `lumos_entitlement_appstore_settings` if the permission had App Store settings:

```terraform
import {
  to = lumos_entitlement.aws_admin
  id = "<entitlement id>"
}

import {
  to = lumos_entitlement_appstore_settings.aws_admin
  id = "<entitlement id>"
}
```

If the configuration set approvals or provisioning (approvers, manager approval, stage 2 approvers, provisioning webhooks, manual provisioning steps), import the workflows Lumos already has, and their bindings. The workflow IDs are on the app (`approval_workflow.id`, `default_provisioning_workflow.id`, `default_deprovisioning_workflow.id`) and on the entitlement (`approval_workflow_override.id`, `provisioning_workflow.id`, `deprovisioning_workflow.id`). An approval template's `branch` step references its chains by ID, and each chain is its own `lumos_workflow`.

```terraform
import {
  to = lumos_workflow.okta_approvals
  id = "<approval template id>"
}

import {
  to = lumos_workflow.okta_manager_then_it
  id = "<approval chain id>"
}

# A binding's import ID is the app's (or entitlement's) ID.
import {
  to = lumos_app_approval_workflow.okta
  id = "<app id>"
}
```

See the [approvals guide](approvals.md) for how templates, chains and bindings fit together.

### 5. Generate the configuration

```shell
terraform plan -generate-config-out=generated.tf
```

Terraform writes a resource block for every import with no matching block. Review `generated.tf`:

* **Connected apps:** `auth` is never read back from Lumos, so the generated `lumos_app` has no credentials. Add `auth`, and `settings` if the integration takes any, from the 0.x configuration.
* **References:** replace hardcoded IDs with references, such as `app_id = lumos_app.okta.id`.
* **Conditions** (`lumos_access_policy.access_condition`, `lumos_pre_approval_rule.condition`, workflow step `inputs`) come back as JSON strings. Rewrite them with `jsonencode(...)` if you like, keeping every attribute and value as `{ label, value }`.

Move the blocks into your own files. Delete the `import` blocks after the first successful apply.

### 6. Plan and apply

```shell
terraform plan
```

Apart from [What to expect on the first apply](#what-to-expect-on-the-first-apply), every resource should show as `will be imported` with no changes. A change on an imported resource means the configuration differs from Lumos: fix the configuration, unless the change is one you want.

```shell
terraform apply
```

## What to expect on the first apply

* **Connected integrations stay connected.** The imported state has no `auth`, so the first apply shows `auth` being added. Terraform records the credentials from the configuration without sending them, so the integration isn't reconnected and no sync starts. An app that's disconnected in Lumos is the exception: its first apply sends the credentials and reconnects it.
* After that, an apply sends credentials only when `auth`, `settings` or `version` changes. A credential change reconnects the integration; for an integration that finishes connecting with browser consent, redo consent in the Lumos UI.

## Resource notes

### `lumos_app`

| 0.x | 1.x |
|---|---|
| `name` | `label_override` (read back as `label`) |
| `description`, `category`, `logo_url`, `website_url` | `description_override`, `category_override`, `logo_url_override`, `website_url_override` (read back without `_override`) |
| `request_instructions` | `lumos_app_appstore_settings.request_instructions` |
| `app_class_id`, `auth`, `settings`, `version` | Same |
| `custom_attributes` (read-only) | Same |
| `links`, `sources`, `user_friendly_label` (read-only) | Removed |
| `allow_multiple_permission_selection` (read-only) | `appstore_settings.allow_multiple_permission_selection` |

* **Overrides work on connected integrations too.** An `*_override` replaces the integration's value in your organization; leaving it out shows the integration's default value. The read-only attributes without `_override` (`label`, `description`, `category`, `logo_url`, `website_url`) always show the effective value.
* **Custom apps** need `label_override`, `category_override` and `description_override`. A custom app can leave `app_class_id` out or set it to `"custom"`, the value it reads back with.
* **Destroy deletes the app** (custom apps) or disconnects it (integrations). In 0.10.x, destroy only removed the app from state. To stop managing an app without deleting it, use a `removed` block.

### `lumos_app_store_app` → `lumos_app_appstore_settings`, `lumos_app` and workflows

| 0.x | 1.x |
|---|---|
| Listing the app (create) / unlisting it (destroy) | `in_app_store` (destroy also unlists) |
| `request_flow.discoverability` (`FULL`, `LIMITED`, `NONE`) | `catalog_visibility` (`visible`, `search-only`, `hidden`) <!-- TODO: confirm the value mapping --> |
| `request_flow.admins` | `admins`, each `{ id, type }` with `type` set to `"user"` or `"group"`. A listed app needs at least one |
| `request_flow.allowed_groups` (`type`, `groups`) | `allowed_requester_group_uuids`. Empty lets anyone request. Groups are referenced by ID only: look one up by its integration ID with `data "lumos_groups"` |
| `custom_request_instructions` | `request_instructions` |
| `request_flow.custom_approval_message` | `lumos_app.post_approval_message` |
| `request_flow.request_validation_inline_webhook.id` | `request_validation_webhook_uuid` |
| `provisioning.allow_multiple_permission_selection` | `allow_multiple_permission_selection` |
| `provisioning.time_based_access` | `time_based_access.options` (labels, as before) |
| `provisioning.groups_provisioning` | See [Groups provisioning](#groups-provisioning) |
| `request_flow.require_manager_approval`, `approvers`, `require_additional_approval`, `approvers_stage_2` | Stages in an approval chain, the fallback of an approval template bound with `lumos_app_approval_workflow` (see [Approvals](#approvals)) |
| `provisioning.provisioning_webhook`, `access_removal_inline_webhook` | `provisioning_webhook` / `deprovisioning_webhook` steps in the app's account workflows |
| `provisioning.manual_steps_needed`, `custom_provisioning_instructions` | A `manual_provisioning_task` step in the app's `account-provisioning` workflow |

`in_app_store`, `catalog_visibility`, `allow_permission_selection` and `allow_multiple_permission_selection` are required. Settings with no 0.x equivalent: `allow_permission_selection`, `redirect_url`, `default_entitlement`, `custom_intake_fields`, `time_based_access.default`.

#### Groups provisioning

0.x's `groups_provisioning` combined two settings that 1.x sets directly:

| 0.x | 1.x |
|---|---|
| `DIRECT_TO_USER` | `visible_in_appstore = false` on every entitlement's `lumos_entitlement_appstore_settings` |
| `GROUPS_AND_HIDDEN` | At least one entitlement with `visible_in_appstore = true`, and `allow_permission_selection = false` on the app |
| `GROUPS_AND_VISIBLE` | At least one entitlement with `visible_in_appstore = true`, and `allow_permission_selection = true` on the app |

#### Approvals

In 0.x, approvals were fields on `lumos_app_store_app`:

```terraform
resource "lumos_app_store_app" "okta" {
  app_id = lumos_app.okta.id
  request_flow = {
    require_manager_approval    = true
    approvers                   = { groups = [{ id = "<security group id>" }] }
    require_additional_approval = true
    approvers_stage_2           = { groups = [{ id = "<IT group id>" }] }
  }
}
```

In 1.x, this configuration maps to an approval template with a single fallback approval chain, itself a workflow. Manager approval maps to a first stage assigned to the `Manager` custom attribute; the 0.x approvers and stage 2 approvers become the stages after it:

```terraform
# The approval template: a workflow that picks which chain a request goes through.
resource "lumos_workflow" "okta_approvals" {
  app_id = lumos_app.okta.id
  type   = "approval"
  name   = "Okta approvals"
  steps = [{
    step = "branch"
    inputs = jsonencode({
      branches             = []
      fallback_workflow_id = lumos_workflow.okta_manager_security_it.id
    })
  }]
}

# The approval chain: a workflow holding the approval stages, in order.
resource "lumos_workflow" "okta_manager_security_it" {
  app_id = lumos_app.okta.id
  type   = "approval-chain"
  name   = "Manager, then Security, then IT"
  steps = [
    {
      step = "approval_task"
      inputs = jsonencode({
        stage_number = 1
        assignees    = [{ label = "Manager", type = "custom_attribute_definition" }]
      })
    },
    {
      step = "approval_task"
      inputs = jsonencode({
        stage_number = 2
        assignees    = [{ uuid = "<security group id>", type = "group" }]
      })
    },
    {
      step = "approval_task"
      inputs = jsonencode({
        stage_number = 3
        assignees    = [{ uuid = "<IT group id>", type = "group" }]
      })
    },
  ]
}

# The binding: makes the template the app's approval process.
resource "lumos_app_approval_workflow" "okta" {
  app_id      = lumos_app.okta.id
  workflow_id = lumos_workflow.okta_approvals.id
}
```

For existing apps, you can import the templates and chains directly from Lumos (see [step 4](#4-write-import-blocks)) and generate their configuration (see [step 5](#5-generate-the-configuration)), rather than writing them by hand.

**New in 1.x: conditional approvals.** 0.x could only manage one approval chain per app. A 1.x template can send different requests through different chains, for example requests from one team through a shorter chain, using `branches`. A template imported from an app whose approvals were set up in the Lumos UI may already have branches: import each chain it references as its own `lumos_workflow`. See the [approvals guide](approvals.md) for how to write branches.

<!-- TODO: confirm the 0.x manager stage reads back through GET /v1/workflows as `{ label = "Manager", type = "custom_attribute_definition" }`. -->

### `lumos_requestable_permission` → `lumos_entitlement` and its settings

| 0.x | 1.x |
|---|---|
| `app_id`, `label` | `lumos_entitlement.app_id`, `label` |
| `request_config.request_approval_config.custom_approval_message` (+ `_override`) | `lumos_entitlement.custom_approval_message_override`. Leave it out to use the app's message |
| `request_config.appstore_visibility` (`VISIBLE`, `HIDDEN`) | `lumos_entitlement_appstore_settings.visible_in_appstore` |
| `request_config.allowed_groups` + `allowed_groups_override` | `allowed_requester_group_uuids_override`. Leave it out to inherit the app's groups; `[]` lets anyone request |
| `request_fulfillment_config.time_based_access` + `time_based_access_override` | `time_based_access_override.options`. Leave it out to inherit the app's options |
| `request_validation_inline_webhook.id` | `request_validation_webhook_uuid` |
| `request_approval_config.*` with `request_approval_config_override` | An approval template bound with `lumos_entitlement_approval_workflow`. Destroying the binding turns the override off |
| `request_fulfillment_config.provisioning_webhook`, `access_removal_inline_webhook`, `manual_steps_needed`, `manual_instructions` | Steps in workflows bound with `lumos_entitlement_provisioning_workflow` / `lumos_entitlement_deprovisioning_workflow` |
| `request_fulfillment_config.provisioning_group` | Import the group's synced entitlement instead of creating one |
| `include_inherited_configs`, `app_class_id`, `app_instance_id`, `type` | Removed. Entitlement settings show only the entitlement's own overrides |

* `lumos_entitlement` creates only manual entitlements. Entitlements synced from an integration can be imported, to manage their settings, but not created or deleted.
* `description_override` is new: a description shown to requesters that survives syncs.
* Destroying a `lumos_entitlement` deletes the entitlement. In 0.x, destroying a permission left its entitlement behind.

### `lumos_access_policy`

| 0.x | 1.x |
|---|---|
| `apps[].permissions[].id` (permission IDs) | `apps[].entitlements[].id` (entitlement IDs) |
| `apps` (list) | `apps` (set; order doesn't matter) |
| `status` (`DRAFT`, `PUBLISHED`), optional | `status` (`draft`, `published`), **required** |
| `is_enabled`, optional | **Required** |
| `access_condition`, `is_everyone_condition` | Same. Set exactly one. Write every attribute and value as `{ label, value }` |

The generated configuration already uses entitlement IDs, so policies don't need [Finding entitlement IDs](#finding-entitlement-ids), except to replace those IDs with references.

### `lumos_pre_approval_rule`

The 0.x fields are replaced by one `condition`. Lumos has already converted existing rules, so the condition generated in step 5 matches what the rule does today.

| 0.x | 1.x |
|---|---|
| `preapproved_groups` | A `condition` term on `identity.groups[*].external_id` with `contains_any`, matching the groups' integration IDs |
| `preapproved_users_by_attribute` | `condition` terms on `identity.custom_attributes["<label>"]` |
| `preapproval_webhooks` | `condition` terms on `webhook.responses["<webhook id>"]` equal to `PREAPPROVED` |
| `time_based_access` | A `condition` term on `access_request.access_length`, in seconds |
| `preapproved_permissions` | `request_includes_entitlements` (entitlement IDs, lowercase) |
| `app_class_id`, `app_instance_id` (read-only) | Removed |

Rules can also pre-approve by on-call schedule in 1.x. As in 0.x, rules created by an access policy are managed through the access policy.

## Finding entitlement IDs

<TODO>

<!-- GOV-4405: likely an `entitlement_uuid` field on 0.x requestable permissions in the API, so each permission ID can be looked up to its entitlement ID. -->

## `lumos_apps` returns every app

In 0.11, `lumos_apps` returned only integrations connected through the API that weren't disconnected (`connection_source = "API"`, `disconnected = false`) unless those arguments were set. In 1.x it has no defaults and returns every app. To keep the 0.11 results, set both:

```terraform
data "lumos_apps" "managed" {
  connection_source = "api"
  disconnected      = false
}
```

`exact_match` is gone; `search` is a case-insensitive substring match against the label, instance identifier and app class ID. To require an exact match, filter the results:

```terraform
locals {
  okta = one([for a in data.lumos_apps.managed.apps : a if a.label == "Okta"])
}
```

## Leaving settings to the Lumos UI

To keep editing some settings in the Lumos UI, leave them out of the configuration and list them in `ignore_changes`:

```terraform
resource "lumos_app_appstore_settings" "okta" {
  # ...
  lifecycle {
    ignore_changes = [request_instructions, redirect_url]
  }
}
```

This works only once the settings are in state, so import them (step 4) rather than creating the resource. See [`lumos_app_appstore_settings`](../resources/app_appstore_settings.md#leaving-settings-to-the-lumos-ui) for details. The same works for `lumos_entitlement_appstore_settings` overrides and `lumos_app`'s `*_override` fields.
