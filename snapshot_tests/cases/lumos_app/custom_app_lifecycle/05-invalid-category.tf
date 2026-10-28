# expect: error
# A category_override that isn't one of the API's categories is refused, and the
# app keeps its current one.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deploy Dashboard v2"
  category_override    = "Not A Category"
  description_override = "Snapshot test custom app"
  logo_url_override    = "https://example.com/deploy.png"
}
