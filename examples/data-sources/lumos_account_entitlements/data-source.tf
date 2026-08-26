data "lumos_account_entitlements" "my_accountentitlements" {
  account_uuid = "...my_account_uuid..."
  expand = [
    "..."
  ]
  page_size  = 100
  page_token = "...my_page_token..."
}