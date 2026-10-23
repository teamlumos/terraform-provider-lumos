---
page_title: "Approvals"
subcategory: ""
description: |-
  Build an app's approval process from an approval template and approval chains.
---

# Approvals

An app's approval process is built from `lumos_workflow` resources in two tiers:

* An **approval template** (`type = "approval"`) decides which chain a request goes through. Its one step is a `branch` step with a list of branches, each pairing a condition with a chain, and a fallback chain. Branches are evaluated in order, and the first one whose condition matches picks the chain. A request that matches no branch goes through the fallback.
* An **approval chain** (`type = "approval-chain"`) holds the approval stages a request goes through.

The template and its chains must all reference the same app.

## The app

The template and its chains belong to one app. Here, a custom app:

```terraform
resource "lumos_app" "deploy_dashboard" {
  label       = "Deploy Dashboard"
  category    = "Developers"
  description = "Internal deploy tool"
}
```

## The approval template

The template's `branch` step references chains by ID. This one sends every request through one chain:

```terraform
resource "lumos_workflow" "deploy_dashboard_approvals" {
  app_id = lumos_app.deploy_dashboard.id
  type   = "approval"
  name   = "Deploy Dashboard approvals"
  steps = [{
    step = "branch"
    inputs = jsonencode({
      branches             = []
      fallback_workflow_id = lumos_workflow.manager_then_it.id
    })
  }]
}
```

* To send some requests through a different chain, add entries to `branches`, each `{ condition = <condition>, workflow_id = <chain id> }`.
* A template with branches must also have a fallback.
* When no approval is required, leave out both: `branches = []` and no `fallback_workflow_id`.

### Branch conditions

Write every attribute and value in a condition as a `{ label, value }` object. Lumos matches on `value`; `label` is the text shown in Lumos, and it's kept as written. A condition written with plain strings, such as `attribute = "identity.team"`, is accepted, but Lumos reads it back in the `{ label, value }` form, so every plan shows a change.

This template sends requests from the Eng and IT teams through an IT-only chain, and every other request through the fallback:

```terraform
resource "lumos_workflow" "deploy_dashboard_approvals" {
  app_id = lumos_app.deploy_dashboard.id
  type   = "approval"
  name   = "Deploy Dashboard approvals"
  steps = [{
    step = "branch"
    inputs = jsonencode({
      branches = [{
        condition = {
          and = [{
            attribute = { label = "Team", value = "identity.team" }
            in        = [{ label = "Engineering", value = "Eng" }, { label = "IT", value = "IT" }]
          }]
        }
        workflow_id = lumos_workflow.it_only.id
      }]
      fallback_workflow_id = lumos_workflow.manager_then_it.id
    })
  }]
}
```

## Approval chains

Each chain lists its approval stages as `approval_task` steps. `stage_number` starts at 1 and counts up with no gaps; Lumos checks this when the template is bound. `assignees` lists who approves that stage:

```terraform
resource "lumos_workflow" "manager_then_it" {
  app_id = lumos_app.deploy_dashboard.id
  type   = "approval-chain"
  name   = "Manager, then IT"
  steps = [
    {
      step = "approval_task"
      inputs = jsonencode({
        stage_number = 1
        assignees    = [{ uuid = "<managers group id>", type = "group" }]
      })
    },
    {
      step = "approval_task"
      inputs = jsonencode({
        stage_number = 2
        assignees    = [{ uuid = "<IT group id>", type = "group" }]
      })
    },
  ]
}
```

## Binding the template

The template does nothing until it's bound. Bind it to the app to make it the approval process for requests to the app:

```terraform
resource "lumos_app_approval_workflow" "deploy_dashboard" {
  app_id      = lumos_app.deploy_dashboard.id
  workflow_id = lumos_workflow.deploy_dashboard_approvals.id
}
```

* An app always has an approval workflow. To change it, set `workflow_id` to a different template; removing the binding from the configuration fails the plan rather than leaving the template bound while Terraform reports it gone. To stop managing it without changing Lumos, use a `removed` block with `lifecycle { destroy = false }`.
* To give one entitlement a different approval process, bind another template with `lumos_entitlement_approval_workflow`. That turns the entitlement's override on; destroying the binding turns it off again, so the entitlement goes back to the app's template.
* A template can be bound to only one app or entitlement.
* Changing a binding's `app_id` fails the plan, including when the new app isn't known until apply, for example an app created or replaced in the same apply. Changing it would replace the binding, and Lumos can't empty the old app's slot, so the old app would keep the template.

### Moving a template to another app

A slot can't be emptied, so moving template `t` from app `a` to app `b` takes two bindings: rebind `a` to a different template, which frees `t`, then bind `t` to `b` with its own binding resource.

```terraform
# App A gets a different template, so T is freed.
resource "lumos_app_approval_workflow" "a" {
  app_id      = lumos_app.a.id
  workflow_id = lumos_workflow.t2.id # was lumos_workflow.t.id
}

# App B gets T, through its own binding resource.
resource "lumos_app_approval_workflow" "b" {
  app_id      = lumos_app.b.id
  workflow_id = lumos_workflow.t.id
  depends_on  = [lumos_app_approval_workflow.a]
}
```

* `depends_on` makes Terraform rebind `a` first. Lumos refuses to bind an approval template that's still bound to another app, so binding `b` first fails.
* Account and entitlement provisioning and deprovisioning workflows can be shared, so their order doesn't matter, but the old app or entitlement keeps the workflow until it's rebound.
* **Replacing the app itself,** for example a `lumos_app` change that forces a new app, fails the plan while a binding points at it. Stop managing the binding first with a `removed` block with `lifecycle { destroy = false }`, replace the app, then add the binding back.

## Limitations

* A chain can be used by only one workflow. Lumos refuses a template that references a chain another workflow already uses.
* A chain a template stops referencing is kept, so it stays managed by its own `lumos_workflow` resource.
