# The same config. Nothing changes, so the lookup reads at plan and the plan
# is empty.
resource "lumos_app" "first" {
  label_override       = "bb-${var.run_id} Lookup One"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "second" {
  label_override       = "bb-${var.run_id} Lookup Two"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_apps" "mine" {
  search = "bb-${var.run_id}"

  depends_on = [lumos_app.first, lumos_app.second]
}

data "lumos_apps" "none" {
  search = "bb-${var.run_id}-no-such-app"
}

output "finds_both_apps" {
  value = toset([for a in data.lumos_apps.mine.apps : a.id]) == toset([lumos_app.first.id, lumos_app.second.id])
}

output "labels" {
  value = sort([for a in data.lumos_apps.mine.apps : a.label])
}

output "no_match_is_empty" {
  value = length(data.lumos_apps.none.apps) == 0
}

output "custom_attributes_are_maps" {
  value = alltrue([for a in data.lumos_apps.mine.apps : can(keys(a.custom_attributes))])
}
