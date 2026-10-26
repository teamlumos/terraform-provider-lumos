# Pre-approves requests for the deploy dashboard's Viewer entitlement from
# anyone on the Platform team, for up to 30 days.
resource "lumos_pre_approval_rule" "platform_viewers" {
  app_id        = lumos_app.deploy_dashboard.id
  justification = "The Platform team reads deploy status every day."

  condition = jsonencode({
    and = [
      {
        attribute = { label = "Team", value = "identity.team" }
        equals    = { label = "Platform", value = "Platform" }
      },
      {
        attribute = { label = "Access length", value = "access_request.access_length" }
        "<="      = { label = "30 days", value = 2592000 }
      },
    ]
  })

  request_includes_entitlements = [lumos_entitlement.deploy_viewer.id]
}

# Pre-approves a request for any of the wiki's entitlements when the
# organization's pre-approval webhook approves it. A request for the wiki
# itself, with no entitlement, isn't covered.
resource "lumos_pre_approval_rule" "wiki_by_webhook" {
  app_id        = lumos_app.wiki.id
  justification = "An access service decides who can read the wiki."

  condition = jsonencode({
    attribute = { label = "Access check", value = "webhook.responses[\"<pre-approval webhook id>\"]" }
    equals    = { label = "Pre-approved", value = "PREAPPROVED" }
  })

  request_has_entitlements = true
}

# Pre-approves requests for the deploy dashboard's Admin entitlement from
# whoever is on call for the platform rotation, for up to 12 hours.
resource "lumos_pre_approval_rule" "deploy_admins_on_call" {
  app_id        = lumos_app.deploy_dashboard.id
  justification = "The on-call engineer may need to roll back a deploy."

  condition = jsonencode({
    and = [
      {
        attribute = { label = "On call", value = "identity.oncall[\"<schedule entitlement id>\"].on_call" }
        equals    = { label = "On call", value = true }
      },
      {
        attribute = { label = "Access length", value = "access_request.access_length" }
        "<="      = { label = "12 hours", value = 43200 }
      },
    ]
  })

  request_includes_entitlements = [lumos_entitlement.deploy_admin.id]
}
