# expect: error
# Publishing the draft would give two published policies the same condition.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Duplicate Condition"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_access_policy" "first" {
  name                   = "bb-${var.run_id} First"
  business_justification = "Snapshot test policy"
  status                 = "published"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{ id = lumos_app.this.id }]
}

resource "lumos_access_policy" "second" {
  name                   = "bb-${var.run_id} Second"
  business_justification = "Snapshot test policy"
  status                 = "published"
  is_enabled             = false

  access_condition = jsonencode({
    attribute = { label = "Email", value = "identity.email" }
    ends_with = { label = "@bb-${var.run_id}.example.com", value = "@bb-${var.run_id}.example.com" }
  })

  apps = [{ id = lumos_app.this.id }]
}
