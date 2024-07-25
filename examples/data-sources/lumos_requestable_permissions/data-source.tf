data "lumos_requestable_permissions" "my_requestablepermissions" {
  app_id       = "...my_app_id..."
  exact_match  = true
  in_app_store = false
  page         = 7
  search_term  = "...my_search_term..."
  size         = 10
}