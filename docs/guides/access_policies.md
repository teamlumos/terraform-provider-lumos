---
page_title: "Access Policies"
---

# Access Policies Guide

This guide explains how to create and manage access policies with the Lumos Terraform provider using the `lumos_access_policy` resource.

## Overview

Access policies define who can access which apps and permissions in Lumos. You configure a policy with a name, the apps and permissions it grants, a business justification, and conditions that restrict which identities the policy applies to.

## Basic Example

The following example creates a policy that applies to everyone (`is_everyone_condition = true`) and does not use custom conditions.

```terraform
resource "lumos_access_policy" "example" {
  name                   = "Engineering Team Policy"
  business_justification = "Default access for engineering"
  is_everyone_condition = true

  apps = [
    {
      id             = "...app_id..."
      is_preapproved = true
      permissions = [
        { id = "...permission_id..." }
      ]
    }
  ]
}
```

## Using Access Conditions

The `access_condition` parameter lets you restrict the policy to specific identities (e.g., by team, title, or custom attributes). Conditions are expressed as a JSON structure.

**You must use Terraform's `jsonencode()` function** to format `access_condition` as a JSON string. The parameter expects a string, not a Terraform object. 

### Example: Policy with Conditions

```terraform
resource "lumos_access_policy" "product_engineering_roles" {
  name                   = "Product and Engineering Roles Policy"
  business_justification = "Grant access to specific roles in Product and Engineering departments"

  access_condition = jsonencode({
    "and" = [
      {
        "attribute" = {
          "label" = "Title"
          "value" = "identity.title"
        }
        "in" = [
          { "label" = "Software Engineer", "value" = "Software Engineer" },
          { "label" = "Product Manager", "value" = "Product Manager" }
        ]
      },
      {
        "attribute" = {
          "label" = "Department"
          "value" = "identity.custom_attributes[\"Department\"]"
        }
        "in" = [
          { "label" = "Product", "value" = "Product" },
          { "label" = "Engineering", "value" = "Engineering" }
        ]
      }
    ]
  })

  apps = [
    {
      id             = "...app_id..."
      is_preapproved = true
      permissions = [{ id = "...permission_id..." }]
    }
  ]
}
```

For full condition syntax and options, see the [Lumos Conditions documentation](https://support.lumos.com/articles/8646284496-building-conditions-in-lumos).
