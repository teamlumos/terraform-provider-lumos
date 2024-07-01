data "lumos_users" "my_users" {
  exact_match = false
  page        = 10
  search_term = "...my_search_term..."
  size        = 6
}