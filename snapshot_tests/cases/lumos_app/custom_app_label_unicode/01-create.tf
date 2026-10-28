# A label with accented letters and an emoji.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Déploiement 🚀"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}
