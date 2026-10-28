# requires: GROUP_ID
# A chain approved by a group, then by the user in the Manager custom
# attribute (found through lumos_custom_attribute_definitions), and a manual
# task assigned through the same attribute.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Assignee Types"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

data "lumos_custom_attribute_definitions" "manager" {
  attribute_type = "user"
  label          = "Manager"
}

locals {
  manager = one(data.lumos_custom_attribute_definitions.manager.custom_attribute_definitions[*].label)
}

resource "lumos_workflow" "chain" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Group, then attribute"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = var.group_id, type = "group" }] })
    },
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 2, assignees = [{ label = local.manager, type = "custom_attribute_definition" }] })
    },
  ]

  lifecycle {
    precondition {
      condition     = local.manager != null
      error_message = "This case needs the user-type custom attribute Manager, as GET /v1/custom-attribute-definitions lists it."
    }
  }
}

resource "lumos_workflow" "manual" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Attribute's task"
  steps = [{
    step = "manual_provisioning_task"
    inputs = jsonencode({
      instructions = "Create the account."
      assignees    = [{ label = local.manager, type = "custom_attribute_definition" }]
    })
  }]

  lifecycle {
    precondition {
      condition     = local.manager != null
      error_message = "This case needs the user-type custom attribute Manager, as GET /v1/custom-attribute-definitions lists it."
    }
  }
}
