# A custom app with two entitlements, found through lumos_app_entitlements
# (all, and by search) and through lumos_entitlements.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Lookup"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "read" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}

resource "lumos_entitlement" "write" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Write"
}

data "lumos_app_entitlements" "all" {
  app_id = lumos_app.this.id

  depends_on = [lumos_entitlement.read, lumos_entitlement.write]
}

data "lumos_app_entitlements" "read" {
  app_id = lumos_app.this.id
  search = "Read"

  depends_on = [lumos_entitlement.read, lumos_entitlement.write]
}

data "lumos_entitlements" "mine" {
  search = "bb-${var.run_id}"

  depends_on = [lumos_entitlement.read, lumos_entitlement.write]
}

output "app_entitlement_labels" {
  value = sort([for e in data.lumos_app_entitlements.all.entitlements : e.label])
}

output "search_finds_read" {
  value = [for e in data.lumos_app_entitlements.read.entitlements : e.id] == [lumos_entitlement.read.id]
}

output "entitlements_finds_both" {
  value = toset([for e in data.lumos_entitlements.mine.entitlements : e.id]) == toset([lumos_entitlement.read.id, lumos_entitlement.write.id])
}
