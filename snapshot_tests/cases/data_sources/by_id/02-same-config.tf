# requires: GROUP_ID
# The same config: the lookups read at plan, and the plan is empty.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} By ID"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id               = lumos_app.this.id
  label                = "bb-${var.run_id} Read"
  description_override = "Read access."
}

data "lumos_current_user" "me" {}

data "lumos_app" "this" {
  id = lumos_app.this.id
}

data "lumos_entitlement" "this" {
  id = lumos_entitlement.this.id
}

data "lumos_user" "me" {
  id = data.lumos_current_user.me.id
}

data "lumos_group" "this" {
  id = var.group_id
}

data "lumos_groups" "this" {
  app_id = data.lumos_group.this.app_id
}

output "app_matches_resource" {
  value = [data.lumos_app.this.label, data.lumos_app.this.category, data.lumos_app.this.description] == [lumos_app.this.label, lumos_app.this.category, lumos_app.this.description]
}

output "entitlement_matches_resource" {
  value = [data.lumos_entitlement.this.label, data.lumos_entitlement.this.app_id, data.lumos_entitlement.this.description] == [lumos_entitlement.this.label, lumos_app.this.id, "Read access."]
}

output "user_matches_current_user" {
  value = data.lumos_user.me.email == data.lumos_current_user.me.email
}

output "group_matches_list" {
  value = contains([for g in data.lumos_groups.this.groups : g.name], data.lumos_group.this.name)
}
