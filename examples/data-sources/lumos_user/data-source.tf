variable "user_id" {
  type = string
}

# Reads one user by their ID. Custom attributes are only read when `expand`
# asks for them.
data "lumos_user" "this" {
  id     = var.user_id
  expand = ["custom_attributes"]
}

output "user_email" {
  value = data.lumos_user.this.email
}
