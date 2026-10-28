# website_url_override is left out after being set, which clears it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deploy Dashboard v2"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
  logo_url_override    = "https://example.com/deploy.png"
}
