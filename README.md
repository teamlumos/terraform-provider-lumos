# Terraform Provider Lumos AppStore (Early Access)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Testing the provider

Note that these instructions allow you to expand on a stub already present in `examples/appstore_permissions`. You can take this example and
create your own `main.tf` file as you see fit!

- Enter into the `/examples/appstore_permissions` directory.
- Fill in an API Token in the `<API TOKEN>` section. Note that you will need a Lumos API Token, which you can fetch using instructions [here](https://developers.lumos.com/reference/overview#creating-an-api-token).

Since our provider is still in early access, we're not on the public Terraform Registry yet. To get up and running, you’ll need to update your provider installer to include an override to build the Lumos provider from a local reference.

- First, find the `GOBIN` path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured.

```shell
$ go env GOBIN
/Users/<Username>/go/bin
```

- If the `GOBIN` go environment variable is not set, use the default path, `/Users/<Username>/go/bin`. You can use the `whoami` command to find your `<Username>`.
- Create a new file called `.terraformrc` in your home directory (`~`), then add to the`dev_overrides` block as shown below. Change the `<PATH>` to the value returned from the `go env GOBIN` command above.

```shell
provider_installation {

  dev_overrides {
      "lumos.com/tf/lumos-appstore" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

You are now ready to test the provider!

```shell
cd examples/appstore_permissions
terraform plan
terraform apply
```

Initially you shouldn't see any planned changes. Once you uncomment the resource blocks and fill in the UUIDs
on the locals block in the starter `main.tf` file you can start testing out the provider.

To get Lumos UUIDs you can call the API:

- Get Group UUIDs: https://developers.lumos.com/reference/getgroups
- Get User UUIDs: https://developers.lumos.com/reference/listusers
- Get AppStore Apps UUIDs: https://developers.lumos.com/reference/getappstoreapps
- Get WebHook UUIDs: https://developers.lumos.com/reference/get_inline_webhooks_inline_webhooks_get

To import existing requestable permissions:

Requirements:

- Python3

1. Get Requestable Permission Stable ID from Lumos API
2. Run scripts/generate_lumos_terraform_resource_for_import.py [Add API Token]

```
generate_lumos_terraform_resource_for_import.py lumos-appstore_requestable_permission <requestable-permission-id>
```

3. Get the output and paste it in examples/appstore_permissions/main.tf
4. Run terraform import

```
terraform import lumos-appstore_requestable_permission.<resource_name> <requestable-permission-id>
```

5. terraform apply
