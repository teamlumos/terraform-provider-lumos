data "lumos_groups" "my_groups" {
  app_id                  = "...my_app_id..."
  exact_match             = false
  integration_specific_id = "...my_integration_specific_id..."
  name                    = "Jamie Nitzsche"
  page                    = 6
  size                    = 2
}