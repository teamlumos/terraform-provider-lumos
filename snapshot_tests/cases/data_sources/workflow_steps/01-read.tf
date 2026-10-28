# Steps available to workflows on a custom app. manual_provisioning_task is
# offered only for account-provisioning.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Workflow Steps"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_workflow_steps" "account_provisioning" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
}

data "lumos_workflow_steps" "account_deprovisioning" {
  app_id = lumos_app.this.id
  type   = "account-deprovisioning"
}

output "provisioning_offers_manual_task" {
  value = contains([for s in data.lumos_workflow_steps.account_provisioning.workflow_steps : s.step], "manual_provisioning_task")
}

output "deprovisioning_does_not_offer_it" {
  value = !contains([for s in data.lumos_workflow_steps.account_deprovisioning.workflow_steps : s.step], "manual_provisioning_task")
}
