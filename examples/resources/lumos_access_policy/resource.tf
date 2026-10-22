# Everyone with an example.com address can request the wiki, pre-approved,
# and the deploy dashboard's Admin entitlement.
resource "lumos_access_policy" "engineering" {
  name                   = "Engineering tools"
  business_justification = "Engineers need the wiki and deploy access to ship."
  status                 = "published"
  is_enabled             = true

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@example.com", value = "@example.com" }
  })

  apps = [
    {
      # No entitlements: grants the app itself.
      id             = lumos_app.wiki.id
      is_preapproved = true
    },
    {
      id           = lumos_app.deploy_dashboard.id
      entitlements = [{ id = lumos_entitlement.deploy_admin.id }]
    },
  ]
}
