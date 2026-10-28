# website_url_override and logo_url_override are added.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deploy Dashboard v2"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
  website_url_override = "https://example.com/deploy"
  logo_url_override    = "https://example.com/deploy.png"
}
