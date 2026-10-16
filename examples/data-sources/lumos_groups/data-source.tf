# Finds a group by its exact name, for example to allow its members to
# request an app.
data "lumos_groups" "engineering" {
  name        = "Engineering"
  exact_match = true
}

output "engineering_group_id" {
  value = one(data.lumos_groups.engineering.groups).id
}
