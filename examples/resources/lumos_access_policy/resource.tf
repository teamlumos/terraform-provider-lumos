resource "lumos_access_policy" "my_accesspolicy" {
  access_condition = {
    map_of_any = {
      key = jsonencode("value")
    }
  }
  apps = [
    {
      id             = "...my_id..."
      is_preapproved = true
      permissions = [
        {
          id = "...my_id..."
        }
      ]
    }
  ]
  business_justification = "...my_business_justification..."
  name                   = "...my_name..."
}