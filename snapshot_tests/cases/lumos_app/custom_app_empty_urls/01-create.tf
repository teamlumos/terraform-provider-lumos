# expect: error
# Empty strings for the optional URLs are refused at plan time; leaving them
# out is how to have none.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Empty URLs"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
  website_url_override = ""
  logo_url_override    = ""
}
