data "lumos_user_accounts" "my_useraccounts" {
  expand = [
    "..."
  ]
  page_size  = 50
  page_token = "...my_page_token..."
  user_id    = "...my_user_id..."
}