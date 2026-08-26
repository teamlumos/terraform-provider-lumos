data "lumos_integrations" "my_integrations" {
  page_size  = 100
  page_token = "...my_page_token..."
  query      = "...my_query..."
}