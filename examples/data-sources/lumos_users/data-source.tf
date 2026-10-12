data "lumos_users" "my_users" {
  exact_match = false
  expand = [
    "..."
  ]
  search_term = "...my_search_term..."
}