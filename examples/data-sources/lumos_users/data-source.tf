data "lumos_users" "my_users" {
  exact_match = false
  expand = [
    "..."
  ]
  page_size   = 50
  page_token  = "...my_page_token..."
  search_term = "...my_search_term..."
}