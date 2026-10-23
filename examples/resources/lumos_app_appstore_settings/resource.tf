variable "engineering_group_id" {
  type = string
}

# Lists Okta in the App Store for one group, with two time-based access
# options and a question on the request form. Time-based access options are
# referred to by the labels your organization defines. Settings left out,
# such as redirect_url, are cleared.
resource "lumos_app_appstore_settings" "okta" {
  app_id                              = lumos_app.okta.id
  in_app_store                        = true
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  request_instructions                = "Request the role your team uses."
  allowed_requester_group_uuids       = [var.engineering_group_id]

  time_based_access = {
    options = ["1 day", "1 week"]
    default = "1 day"
  }

  custom_intake_fields = [{
    reference_key = "justification"
    field_type    = "text"
    prompt        = "Why do you need access?"
    is_required   = true
  }]
}
