data "lumos_app_entitlement_app_store_catalog" "my_appentitlementappstorecatalog" {
  app_uuid         = "...my_app_uuid..."
  page_size        = 100
  page_token       = "...my_page_token..."
  target_user_uuid = "...my_target_user_uuid..."
}