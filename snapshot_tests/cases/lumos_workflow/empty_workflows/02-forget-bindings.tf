# The bindings are forgotten so the case can be destroyed.
# Lumos refuses to delete a workflow while it's bound, so the bound workflows are
# forgotten too, and they stay in Lumos.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Empty Workflows"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

removed {
  from = lumos_app_provisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_app_deprovisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_entitlement_provisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_workflow.account_provisioning

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_workflow.account_deprovisioning

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_workflow.entitlement_provisioning

  lifecycle {
    destroy = false
  }
}
