# A dropdown field is added after it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Intake Fields"
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
      reference_key    = "level"
      field_type       = "dropdown"
      prompt           = "Which access level?"
      description      = "Pick the lowest level that works."
      dropdown_options = ["Read", "Write"]
      is_required      = false
    },
  ]
}

data "lumos_current_user" "me" {}
