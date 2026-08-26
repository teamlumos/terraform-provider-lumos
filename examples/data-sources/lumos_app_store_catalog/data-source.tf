data "lumos_app_store_catalog" "my_appstorecatalog" {
  page_size        = 100
  page_token       = "...my_page_token..."
  target_user_uuid = "...my_target_user_uuid..."
}