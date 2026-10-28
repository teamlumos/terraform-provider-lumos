# URLs with an uppercase host, a trailing slash and a query string. If the API
# normalizes them, the apply fails with an inconsistent result.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} URL Forms"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
  website_url_override = "https://Example.com/deploy/?ref=terraform"
  logo_url_override    = "https://EXAMPLE.com/deploy.png?v=1"
}
