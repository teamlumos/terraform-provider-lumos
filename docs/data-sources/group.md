---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lumos_group Data Source - terraform-provider-lumos"
subcategory: ""
description: |-
  Group DataSource
---

# lumos_group (Data Source)

Group DataSource

## Example Usage

```terraform
data "lumos_group" "my_group" {
  group_id = "...my_group_id..."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String)

### Read-Only

- `app_id` (String) The ID of the app that sources this group.
- `description` (String) The description of this group.
- `group_lifecycle` (String) The lifecycle of this group.
- `id` (String) The ID of this group.
- `integration_specific_id` (String) The ID of this group, specific to the integration.
- `name` (String) The name of this group.
- `source_app_id` (String) The ID of the app that sources this group.
