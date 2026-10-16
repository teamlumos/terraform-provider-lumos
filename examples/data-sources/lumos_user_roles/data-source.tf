data "lumos_current_user" "me" {}

# The Lumos roles assigned to a user.
data "lumos_user_roles" "me" {
  user_id = data.lumos_current_user.me.id
}

output "my_roles" {
  value = data.lumos_user_roles.me.roles
}
