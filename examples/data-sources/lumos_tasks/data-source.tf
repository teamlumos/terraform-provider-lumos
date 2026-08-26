data "lumos_tasks" "my_tasks" {
  access_request_id = "...my_access_request_id..."
  assigned_group_ids = [
    "..."
  ]
  assigned_user_ids = [
    "..."
  ]
  category = "...my_category..."
  expand = [
    "..."
  ]
  order_by       = "...my_order_by..."
  page_size      = 100
  page_token     = "...my_page_token..."
  product_area   = "...my_product_area..."
  search         = "...my_search..."
  show_all_tasks = false
  status         = "...my_status..."
  type           = "...my_type..."
}