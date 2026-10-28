# Only description_override changes.
resource "lumos_app" "this" {
  label_override        = "bb-${var.run_id} Edits"
  category_override     = "IT & Security"
  description_override  = "Snapshot test custom app, edited"
  post_approval_message = "Your access is ready."
  website_url_override  = "https://example.com/deploy"
  logo_url_override     = "https://example.com/deploy.png"
}
