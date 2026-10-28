# expect: error
# An empty description is refused at plan time.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Empty Description"
  category_override    = "Developers"
  description_override = ""
}
