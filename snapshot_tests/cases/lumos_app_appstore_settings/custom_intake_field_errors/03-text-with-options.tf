# expect: error
# A text field can't have options.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Intake Field Errors"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app_appstore_settings" "this" {
  app_id                              = lumos_app.this.id
  in_app_store                        = true
  admins                              = [{ id = data.lumos_current_user.me.id, type = "user" }]
  catalog_visibility                  = "visible"
  allow_permission_selection          = true
  allow_multiple_permission_selection = false
  custom_intake_fields = [
    {
      reference_key = "team"
      field_type    = "text"
      prompt        = "Which team are you on?"
      description   = "Used to route the request."
      is_required   = true
    },
    {
      reference_key    = "note"
      field_type       = "text"
      prompt           = "Anything else?"
      dropdown_options = ["Yes", "No"]
      is_required      = false
    },
  ]
}

data "lumos_current_user" "me" {}
