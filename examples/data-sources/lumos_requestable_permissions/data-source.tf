data "lumos_requestable_permissions" "my_requestablepermissions" {
  app_id       = "...my_app_id..."
  exact_match  = false
  in_app_store = false
  page         = 9
  search_term  = "...my_search_term..."
  size         = 46
}