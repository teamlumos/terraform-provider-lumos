resource "lumos_access_policy" "my_accesspolicy" {
  access_condition = {
    key = jsonencode("value")
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
  is_everyone_condition  = false
  name                   = "...my_name..."
}