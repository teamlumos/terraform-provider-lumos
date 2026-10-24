variable "group_id" {
  type = string
}

# Reads one group by its ID, for example to show its name next to the ID a
# config already holds.
data "lumos_group" "this" {
  id = var.group_id
}

output "group_name" {
  value = data.lumos_group.this.name
}
