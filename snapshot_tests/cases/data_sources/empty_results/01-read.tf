# Filters that match nothing return empty lists, not errors. The accounts
# lookup uses a custom app the case creates, which has no accounts.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} No Accounts"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_accounts" "none" {
  app_id = lumos_app.this.id
}

data "lumos_users" "none" {
  search = "bb-${var.run_id}-nobody"
}

data "lumos_groups" "none" {
  name        = "bb-${var.run_id}-no-such-group"
  exact_match = true
}

data "lumos_access_policies" "none" {
  name = "bb-${var.run_id}-no-such-policy"
}

output "accounts_empty" {
  value = length(data.lumos_accounts.none.accounts) == 0
}

output "users_empty" {
  value = length(data.lumos_users.none.users) == 0
}

output "groups_empty" {
  value = length(data.lumos_groups.none.groups) == 0
}

output "access_policies_empty" {
  value = length(data.lumos_access_policies.none.access_policies) == 0
}
