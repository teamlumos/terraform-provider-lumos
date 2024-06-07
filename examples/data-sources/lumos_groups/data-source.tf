data "lumos_groups" "my_groups" {
  app_id                  = "...my_app_id..."
  exact_match             = true
  integration_specific_id = "...my_integration_specific_id..."
  name                    = "Mrs. Wilson Krajcik DDS"
  page                    = 1
  size                    = 4
}