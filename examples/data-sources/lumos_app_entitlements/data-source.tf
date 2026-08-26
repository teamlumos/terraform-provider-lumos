data "lumos_app_entitlements" "my_appentitlements" {
  app_uuid = "...my_app_uuid..."
  expand = [
    "..."
  ]
  page_size           = 100
  page_token          = "...my_page_token..."
  visible_in_appstore = true
}