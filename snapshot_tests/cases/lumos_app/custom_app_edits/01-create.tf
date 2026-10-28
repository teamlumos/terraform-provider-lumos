# A custom app with every field a custom app can set.
resource "lumos_app" "this" {
  label_override        = "bb-${var.run_id} Edits"
  category_override     = "Developers"
  description_override  = "Snapshot test custom app"
  post_approval_message = "Your access is ready."
  website_url_override  = "https://example.com/deploy"
  logo_url_override     = "https://example.com/deploy.png"
}
