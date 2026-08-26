data "lumos_workflows" "my_workflows" {
  app_id     = "...my_app_id..."
  page_size  = 100
  page_token = "...my_page_token..."
  type       = "...my_type..."
}