---
page_title: "Access Policy Conditions"
---

# Access Policy Conditions Guide

The `access_condition` parameter in `lumos_access_policy` resources allows you to define complex rules for filtering and matching identities. This document provides detailed examples of how to use `access_condition` with Terraform.

## Important: Using `jsonencode()`

**You must use Terraform's `jsonencode()` function** to properly format the `access_condition` as a JSON string. The `access_condition` parameter expects a JSON string, not a Terraform object. Using `jsonencode()` ensures that your condition is properly serialized.

## Example: Access Policy with Multiple Conditions

```terraform
resource "lumos_access_policy" "my_accesspolicy" {
  access_condition = jsonencode({
    "and" = [
      {
        "attribute" = {
          "label" = "Title"
          "value" = "identity.title"
        }
        "in" = [
          {
            "label" = "Software Engineer"
            "value" = "Software Engineer"
          },
          {
            "label" = "Product Manager"
            "value" = "Product Manager"
          }
        ]
      },
      {
        "attribute" = {
          "label" = "Department"
          "value" = "identity.custom_attributes[\"Department\"]"
        }
        "in" = [
          {
            "label" = "Product"
            "value" = "Product"
          },
          {
            "label" = "Engineering"
            "value" = "Engineering"
          }
        ]
      }
    ]
  })
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
  business_justification = "Grant access to specific roles in Product and Engineering departments"
  name                   = "Product and Engineering Roles Policy"
}
```

## Understanding the Condition Structure

For more information, see the [Lumos Conditions documentation](https://support.lumos.com/articles/8646284496-building-conditions-in-lumos).
