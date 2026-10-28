# A rule read back through lumos_pre_approval_rule (by ID) and
# lumos_pre_approval_rules (the app's rules). On create the outputs are known
# only after apply.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Pre-Approval Data Sources"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "viewer" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Viewer"
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_pre_approval_rule" "this" {
  app_id        = lumos_app.this.id
  justification = "Snapshot test rule"
  condition = jsonencode({
    and = [
      {
        attribute = { label = "Email", value = "identity.email" }
        equals    = { label = "bb-${var.run_id}@example.com", value = "bb-${var.run_id}@example.com" }
      },
      {
        attribute = { label = "Access length", value = "access_request.access_length" }
        "<="      = { label = "30 days", value = 2592000 }
      },
    ]
  })
  request_includes_entitlements = [lumos_entitlement.viewer.id]
}

data "lumos_pre_approval_rule" "this" {
  app_id = lumos_pre_approval_rule.this.app_id
  id     = lumos_pre_approval_rule.this.id
}

data "lumos_pre_approval_rules" "this" {
  app_id = lumos_pre_approval_rule.this.app_id
}

output "by_id_matches_resource" {
  value = [data.lumos_pre_approval_rule.this.justification, jsondecode(data.lumos_pre_approval_rule.this.condition), data.lumos_pre_approval_rule.this.request_includes_entitlements] == [lumos_pre_approval_rule.this.justification, jsondecode(lumos_pre_approval_rule.this.condition), lumos_pre_approval_rule.this.request_includes_entitlements]
}

output "list_finds_only_this_rule" {
  value = [for r in data.lumos_pre_approval_rules.this.pre_approval_rules : r.id] == [lumos_pre_approval_rule.this.id]
}
