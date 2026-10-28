# The same config. Nothing changes, so the lookups read at plan and the plan
# is empty.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Policy Lookup"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

resource "lumos_access_policy" "condition" {
  name                   = "bb-${var.run_id} Lookup Condition"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{
    id             = lumos_app.this.id
    is_preapproved = true
    entitlements   = [{ id = lumos_entitlement.this.id }]
  }]
}

resource "lumos_access_policy" "everyone" {
  name                   = "bb-${var.run_id} Lookup Everyone"
  business_justification = "Snapshot test policy"
  status                 = "draft"
  is_enabled             = false
  is_everyone_condition  = true

  apps = [{ id = lumos_app.this.id }]
}

data "lumos_access_policies" "mine" {
  name = "bb-${var.run_id} Lookup"

  depends_on = [lumos_access_policy.condition, lumos_access_policy.everyone]
}

data "lumos_access_policies" "upper_case" {
  name = upper("bb-${var.run_id} Lookup Condition")

  depends_on = [lumos_access_policy.condition]
}

data "lumos_access_policies" "none" {
  name = "bb-${var.run_id} No Such Policy"
}

data "lumos_access_policy" "condition" {
  id = lumos_access_policy.condition.id
}

data "lumos_access_policy" "everyone" {
  id = lumos_access_policy.everyone.id
}

output "name_filter_finds_both" {
  value = toset([for p in data.lumos_access_policies.mine.access_policies : p.id]) == toset([lumos_access_policy.condition.id, lumos_access_policy.everyone.id])
}

output "name_filter_ignores_case" {
  value = [for p in data.lumos_access_policies.upper_case.access_policies : p.id] == [lumos_access_policy.condition.id]
}

output "no_match_is_empty" {
  value = length(data.lumos_access_policies.none.access_policies) == 0
}

output "by_id_matches_resource" {
  value = {
    name                   = data.lumos_access_policy.condition.name == lumos_access_policy.condition.name
    business_justification = data.lumos_access_policy.condition.business_justification == lumos_access_policy.condition.business_justification
    status                 = data.lumos_access_policy.condition.status == lumos_access_policy.condition.status
    is_enabled             = data.lumos_access_policy.condition.is_enabled == lumos_access_policy.condition.is_enabled
    is_everyone_condition  = data.lumos_access_policy.condition.is_everyone_condition == lumos_access_policy.condition.is_everyone_condition
    access_condition       = jsondecode(data.lumos_access_policy.condition.access_condition) == jsondecode(lumos_access_policy.condition.access_condition)
    apps                   = jsonencode(data.lumos_access_policy.condition.apps) == jsonencode(lumos_access_policy.condition.apps)
  }
}

output "everyone_reads_back" {
  value = {
    is_everyone_condition = data.lumos_access_policy.everyone.is_everyone_condition
    access_condition      = data.lumos_access_policy.everyone.access_condition
  }
}
