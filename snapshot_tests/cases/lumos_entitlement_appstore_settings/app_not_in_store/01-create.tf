# expect: error
# Entitlement settings for an app that isn't in the App Store: the API returns
# 404.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Unlisted App"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

resource "lumos_entitlement_appstore_settings" "this" {
  entitlement_id      = lumos_entitlement.this.id
  visible_in_appstore = true
}
