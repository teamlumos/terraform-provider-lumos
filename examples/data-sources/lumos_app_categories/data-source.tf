# The categories a custom app can be given.
data "lumos_app_categories" "all" {}

output "category_labels" {
  value = [for c in data.lumos_app_categories.all.app_categories : c.label]
}
