data "lumos_groups" "my_groups" {
  app_id                  = "...my_app_id..."
  exact_match             = false
  integration_specific_id = "...my_integration_specific_id..."
  name                    = "...my_name..."
  page                    = 1
  size                    = 2
}