variable "app_id" {
  type = string
}

# The provisioning steps available for an app.
data "lumos_workflow_steps" "provisioning" {
  app_id = var.app_id
  type   = "account-provisioning"
}

output "provisioning_steps" {
  value = [for s in data.lumos_workflow_steps.provisioning.workflow_steps : s.step]
}
