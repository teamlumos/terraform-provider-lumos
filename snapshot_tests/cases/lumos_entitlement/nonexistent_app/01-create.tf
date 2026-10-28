# expect: error
# An entitlement on an app that doesn't exist.
resource "lumos_entitlement" "this" {
  app_id = "00000000-0000-0000-0000-000000000000"
  label  = "bb-${var.run_id} Orphan"
}
