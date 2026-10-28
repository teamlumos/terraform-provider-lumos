# A description with several lines.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Multiline"
  category_override    = "Developers"
  description_override = <<-EOT
    Snapshot test custom app.

    Second paragraph, after a blank line.
  EOT
}
