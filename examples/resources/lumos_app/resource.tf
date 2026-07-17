resource "lumos_app" "my_app" {
  app_class_id         = "...my_app_class_id..."
  auth                 = "{ \"see\": \"documentation\" }"
  category             = "...my_category..."
  description          = "...my_description..."
  logo_url             = "...my_logo_url..."
  name                 = "...my_name..."
  request_instructions = "...my_request_instructions..."
  settings             = "{ \"see\": \"documentation\" }"
  version              = "...my_version..."
  website_url          = "...my_website_url..."
}