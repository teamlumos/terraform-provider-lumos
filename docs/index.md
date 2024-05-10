# Lumos AppStore Early Access Provider

## Example Usage

```terraform
terraform {
  required_providers {
    lumos-appstore = {
      source = "lumos.com/tf/lumos-appstore"
    }
  }
}

provider "lumos-appstore" {
  api_token = "<LUMOS API TOKEN>"
  base_url = "https://api.lumos.com"
}
```

## Schema

- `api_token` (Required) <Lumos API Token>
- `base_url` (Optional) "https://api.lumos.com"
