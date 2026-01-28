# Lumos Terraform Provider

<div align="left">
  <a href="https://speakeasyapi.dev/">
    <img src="https://custom-icon-badges.demolab.com/badge/-Built%20By%20Speakeasy-212015?style=for-the-badge&logoColor=FBE331&logo=speakeasy&labelColor=545454" />
  </a>
</div>

The **Lumos Terraform Provider** allows administrators and IAM engineers to manage Lumos configuration—such as apps, permissions, access policies, and pre-approval rules—using Infrastructure as Code.

This provider is **generated from Lumos’s OpenAPI specification** using [Speakeasy](https://speakeasyapi.dev/) and is published publicly to the Terraform Registry.

---

## Table of Contents

- [Generating the Provider](#generating-the-provider)
- [Installation](#installation)
- [Authentication](#authentication)
- [Available Resources and Data Sources](#available-resources-and-data-sources)
- [Release](#release)
- [Local Development & Testing](#local-development--testing)
- [Contributions](#contributions)
- [SDK Created by Speakeasy](#sdk-created-by-speakeasy)

---

## Generating the Provider

This provider is generated from the Lumos OpenAPI specification using Speakeasy.
The OpenAPI JSON is generated from the FastAPI-based `apiserver` in the main Lumos repo.

### Regenerating the Terraform Provider

1. Copy the updated OpenAPI JSON into this repo as `original-openapi.json`
2. Update `overlay.yaml` if needed to opt in new resources, data sources, hooks, or validators
3. Run:

```sh
make provider
```

This command:

* Runs Speakeasy using `.speakeasy/workflow.yaml`
* Applies `overlay.yaml` to `original-openapi.json`
* Regenerates provider code into `internal/`
* Runs `go mod tidy`
* Builds the provider
* Runs tests
* Generates documentation in `docs/`

### Working with the Terraform overlay

The Terraform overlay (`overlay.yaml`) defines how the Lumos OpenAPI specification is mapped into Terraform concepts. This includes:

- Which API endpoints are exposed as Terraform resources or data sources
- Schema and naming overrides
- Custom hooks and validators
- Terraform-specific behavior that cannot be expressed directly in OpenAPI

#### `make overlay`

To preview changes to the overlay, first place an OpenAPI spec in the repo as `openapi.json`, then run:

```sh
make overlay
```

This command combines `openapi.json` with `overlay.yaml` and generates a diff showing how the provider output would change. It does not regenerate files or modify the working tree.

> Note: `make overlay` expects `openapi.json`, while full provider generation uses `original-openapi.json`.

---

## Installation

To install the published provider, add the following to your Terraform configuration and run `terraform init`:

```hcl
terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.11.2"
    }
  }
}

provider "lumos" {
  # Optional - defaults to https://api.lumos.com
}
```

---

## Authentication

Authentication can be configured via environment variables or provider configuration.

Configuration precedence:

1. Provider configuration
2. Environment variables

| Method                             | Configuration                     |
| ---------------------------------- | --------------------------------- |
| Environment variable (recommended) | `export LUMOS_ACCESS_TOKEN="..."` |
| Provider configuration             | `http_bearer`                     |

---

## Available Resources and Data Sources

### Resources

* [lumos_access_policy](docs/resources/access_policy.md)
* [lumos_app](docs/resources/app.md)
* [lumos_app_store_app](docs/resources/app_store_app.md)
* [lumos_pre_approval_rule](docs/resources/pre_approval_rule.md)
* [lumos_requestable_permission](docs/resources/requestable_permission.md)

### Data Sources

* [lumos_access_policies](docs/data-sources/access_policies.md)
* [lumos_access_policy](docs/data-sources/access_policy.md)
* [lumos_app](docs/data-sources/app.md)
* [lumos_apps](docs/data-sources/apps.md)
* [lumos_app_store_app](docs/data-sources/app_store_app.md)
* [lumos_app_store_app_settings](docs/data-sources/app_store_app_settings.md)
* [lumos_group](docs/data-sources/group.md)
* [lumos_groups](docs/data-sources/groups.md)
* [lumos_requestable_permission](docs/data-sources/requestable_permission.md)
* [lumos_requestable_permissions](docs/data-sources/requestable_permissions.md)
* [lumos_user](docs/data-sources/user.md)
* [lumos_users](docs/data-sources/users.md)

---

## Release

Releases are currently performed manually.

1. Create a new Git tag
2. Push the tag to GitHub

This triggers the
[`terraform_publish.yaml`](https://github.com/teamlumos/terraform-provider-lumos/actions/workflows/terraform_publish.yaml)
GitHub Action, which publishes the provider to the Terraform Registry.

---

## Local Development & Testing

This section describes how to build and test the Lumos Terraform Provider **locally**.

### Step 1 (optional): Build and use a local provider binary

This step is only required if testing against a locally built Lumos Terraform Provider.
If you skip this step, Terraform will use the **published provider from the Terraform Registry**.

#### Build the provider

From the root of the provider repo, run:

```sh
go build
```

This builds a local `terraform-provider-lumos` binary.

#### Configure Terraform to use the local build

By default, Terraform downloads the Lumos provider from the registry.
To force Terraform to use your **locally built provider binary**, configure a `dev_overrides` entry in your local Terraform CLI config file (`.terraformrc`).

Terraform looks for this file in your **home directory**:

* macOS / Linux: `~/.terraformrc`
* Windows: `%APPDATA%\terraform.rc`

If the file does not exist, create it and add:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/teamlumos/lumos" = "/Users/johndoe/Desktop/terraform-provider-lumos"
  }

  direct {}
}
```

---

### Step 2: Set your access token

Generate an access token for your intended environment by going to **Settings → API Tokens** in the Lumos UI.

Then export it in your shell:

```sh
export LUMOS_ACCESS_TOKEN="MY_ACCESS_TOKEN"
```

---

### Step 3: Configure Terraform files in `sandbox/`

All local Terraform testing should be done from the `sandbox/` directory.

* `sandbox/` is **gitignored**
* It is safe to modify freely

Your `sandbox/` directory must include a `provider.tf` file and any resources or data sources you want to test.

#### Example `provider.tf`

```hcl
terraform {
  required_providers {
    lumos = {
      source = "teamlumos/lumos"
    }
  }
}

provider "lumos" {
  # Uses https://api.lumos.com by default
}
```

> **Internal (Lumos employees only):**
> To test against a local Lumos API server, run the apiserver from the main Lumos repo and set:
>
> ```hcl
> provider "lumos" {
>   server_url = "http://localhost:8000"
> }
> ```
>
> External users should not configure a localhost `server_url`.

#### Example resource file (`resource.tf`)

```hcl
resource "lumos_access_policy" "example" {
  name                   = "example-policy"
  business_justification = "Test"
  # other fields omitted for brevity
}
```

---

### Step 4: Run Terraform from `sandbox/`

From the provider repo root:

```sh
cd sandbox
terraform init
terraform plan
terraform apply
```

> Note: You can skip `terraform init` if Step 1 was performed (i.e. dev_overrides is enabled).

---

## Contributions

While we welcome contributions, this provider is largely generated.

Feel free to open a PR or GitHub issue as a proof of concept, and we’ll work to incorporate changes into future releases.

---

## SDK Created by Speakeasy

This SDK is generated and maintained using [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks).
<!-- Start Summary [summary] -->
## Summary

Lumos: The Lumos provider allows you to manage resources such as Apps, Permissions, and Pre-Approval Rules
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [Lumos Terraform Provider](#lumos-terraform-provider)
  * [Generating the Provider](#generating-the-provider)
  * [Installation](#installation)
  * [Authentication](#authentication)
  * [Available Resources and Data Sources](#available-resources-and-data-sources)
  * [Release](#release)
  * [Local Development & Testing](#local-development-testing)
* [Copy the TF_REATTACH_PROVIDERS env var](#copy-the-tfreattachproviders-env-var)
* [In a new terminal](#in-a-new-terminal)

<!-- End Table of Contents [toc] -->

<!-- Start Installation [installation] -->
## Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.9.0"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}
```
<!-- End Installation [installation] -->

<!-- Start Authentication [security] -->
## Authentication

This provider supports authentication configuration via environment variables and provider configuration.

The configuration precedence is:

- Provider configuration
- Environment variables

Available configuration:

| Provider Attribute | Description |
|---|---|
| `http_bearer` | HTTP Bearer. Configurable via environment variable `LUMOS_ACCESS_TOKEN`. |
<!-- End Authentication [security] -->

<!-- Start Available Resources and Data Sources [operations] -->
## Available Resources and Data Sources

### Resources

* [lumos_access_policy](docs/resources/access_policy.md)
* [lumos_app](docs/resources/app.md)
* [lumos_app_store_app](docs/resources/app_store_app.md)
* [lumos_pre_approval_rule](docs/resources/pre_approval_rule.md)
* [lumos_requestable_permission](docs/resources/requestable_permission.md)
### Data Sources

* [lumos_access_policies](docs/data-sources/access_policies.md)
* [lumos_access_policy](docs/data-sources/access_policy.md)
* [lumos_app](docs/data-sources/app.md)
* [lumos_apps](docs/data-sources/apps.md)
* [lumos_app_store_app](docs/data-sources/app_store_app.md)
* [lumos_app_store_app_settings](docs/data-sources/app_store_app_settings.md)
* [lumos_group](docs/data-sources/group.md)
* [lumos_groups](docs/data-sources/groups.md)
* [lumos_requestable_permission](docs/data-sources/requestable_permission.md)
* [lumos_requestable_permissions](docs/data-sources/requestable_permissions.md)
* [lumos_user](docs/data-sources/user.md)
* [lumos_users](docs/data-sources/users.md)
<!-- End Available Resources and Data Sources [operations] -->

<!-- Start Testing the provider locally [usage] -->
## Testing the provider locally

#### Local Provider

Should you want to validate a change locally, the `--debug` flag allows you to execute the provider against a terraform instance locally.

This also allows for debuggers (e.g. delve) to be attached to the provider.

```sh
go run main.go --debug
# Copy the TF_REATTACH_PROVIDERS env var
# In a new terminal
cd examples/your-example
TF_REATTACH_PROVIDERS=... terraform init
TF_REATTACH_PROVIDERS=... terraform apply
```

#### Compiled Provider

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

1. Execute `go build` to construct a binary called `terraform-provider-lumos`
2. Ensure that the `.terraformrc` file is configured with a `dev_overrides` section such that your local copy of terraform can see the provider binary

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/teamlumos/lumos" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
<!-- End Testing the provider locally [usage] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->
