# Searches entitlements across every app. Keep the search specific: one that
# takes the API too long fails, and lumos_app_entitlements is faster for one
# app.
data "lumos_entitlements" "admins" {
  search = "Admin"
}

output "admin_entitlement_ids" {
  value = [for e in data.lumos_entitlements.admins.entitlements : e.id]
}
