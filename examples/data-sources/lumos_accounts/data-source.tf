variable "app_id" {
  type = string
}

# A user's accounts in one app, found by email.
data "lumos_accounts" "alice" {
  app_id = var.app_id
  email  = "alice@example.com"
}

output "alice_account_ids" {
  value = [for a in data.lumos_accounts.alice.accounts : a.id]
}
