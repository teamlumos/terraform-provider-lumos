# Policies whose name matches, ignoring case. Leave name out to list every
# access policy.
data "lumos_access_policies" "engineering" {
  name = "Engineering tools"
}
