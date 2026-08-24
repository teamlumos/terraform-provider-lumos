resource "lumos_app" "my_app" {
  app_id                  = "...my_app_id..."
  app_instance_identifier = ""
  auth = {
    # ...
  }
  settings = {
    # ...
  }
  version = "...my_version..."
}