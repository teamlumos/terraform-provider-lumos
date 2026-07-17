data "lumos_apps" "my_apps" {
  connection_source = "API"
  disconnected      = false
  exact_match       = false
  expand = [
    "..."
  ]
  name_search = "...my_name_search..."
  page        = 3
  size        = 19
}