data "lumos_apps" "my_apps" {
  connection_source = "API"
  disconnected      = false
  expand = [
    "..."
  ]
  in_app_store = true
  page_size    = 100
  page_token   = "...my_page_token..."
}