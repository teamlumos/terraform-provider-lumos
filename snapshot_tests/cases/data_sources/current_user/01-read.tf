# The token's user, found again by email through lumos_users, with their roles.
data "lumos_current_user" "me" {}

data "lumos_users" "me" {
  search      = data.lumos_current_user.me.email
  exact_match = true
}

data "lumos_user_roles" "me" {
  user_id = data.lumos_current_user.me.id
}

output "email_looks_like_an_email" {
  value = can(regex("^[^@]+@[^@]+$", data.lumos_current_user.me.email))
}

output "users_finds_exactly_me" {
  value = [for u in data.lumos_users.me.users : u.id] == [data.lumos_current_user.me.id]
}

output "has_roles" {
  value = length(data.lumos_user_roles.me.roles) > 0
}
