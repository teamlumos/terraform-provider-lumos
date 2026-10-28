# post_approval_message is left out after being set, which removes it.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Edits"
  category_override    = "IT & Security"
  description_override = "Snapshot test custom app, edited"
  website_url_override = "https://example.com/deploy-v2"
}
