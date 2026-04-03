data "lumos_users" "my_users" {
  exact_match = false
  expand = [
    "..."
  ]
  page        = 1
  search_term = "...my_search_term..."
  size        = 50
}