data "lumos_custom_attribute_definitions" "my_customattributedefinitions" {
  app_id         = "87b49206-e897-4b47-b1e1-7766eb3ad241"
  attribute_type = "GROUP"
  config_type    = "FOUND_ACCOUNT"
  label          = "...my_label..."
  page_size      = 100
  page_token     = "...my_page_token..."
  source_type    = "NATIVE"
}