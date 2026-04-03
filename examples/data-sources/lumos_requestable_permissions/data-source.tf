data "lumos_requestable_permissions" "my_requestablepermissions" {
  app_id                    = "...my_app_id..."
  exact_match               = false
  in_app_store              = false
  include_inherited_configs = true
  page                      = 1
  search_term               = "...my_search_term..."
  size                      = 50
}