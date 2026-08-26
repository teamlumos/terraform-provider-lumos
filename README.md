# Lumos Terraform Provider

> ## ⚠️ ALPHA — 1.x targets the Lumos V1 API
>
> The 1.x line may introduce breaking changes without prior notice. Pin an
> exact version: `version = "1.0.0-alpha.1"`.
>
> **For production, stay on 0.10.x:** `version = "~> 0.10"`
>
> The 1.x line exposes **only** V1 API resources — it covers the whole `/v1`
> surface, but the 0.10.x resources built on the un-prefixed endpoints (access
> policy, pre-approval rule, requestable permission, group) are not carried over.
> A single configuration cannot use 0.10.x and 1.x together — Terraform loads one
> version per provider address. Use separate root modules if you need both.
>
> The V1 API is itself gated: endpoints return 404 unless the
> `external-v1-api-enabled` flag is on for your tenant.

<div align="left">
  <a href="https://speakeasyapi.dev/">
    <img src="https://custom-icon-badges.demolab.com/badge/-Built%20By%20Speakeasy-212015?style=for-the-badge&logoColor=FBE331&logo=speakeasy&labelColor=545454" />
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge" />
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
- [License](#license)
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
      version = "1.0.0-alpha.1"
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

Generated from the V1 API surface: **9 managed resources** and **33 data
sources**. The authoritative list is generated further down this file, under
[Available Resources and Data Sources](#available-resources-and-data-sources-1)
— it is rewritten on every `make provider`, so it never drifts.

Three resources are full-lifecycle V1 entities: `lumos_app`,
`lumos_entitlement`, and `lumos_workflow`.

The other six bind a workflow to a slot on an app or an entitlement:

| Resource | Slot |
| --- | --- |
| `lumos_app_approval_workflow` | app approval |
| `lumos_app_provisioning_workflow` | app default provisioning |
| `lumos_app_deprovisioning_workflow` | app default deprovisioning |
| `lumos_entitlement_approval_workflow` | entitlement approval override |
| `lumos_entitlement_provisioning_workflow` | entitlement provisioning |
| `lumos_entitlement_deprovisioning_workflow` | entitlement deprovisioning |

A binding points at a workflow you created separately — set `workflow_id` to an
existing `lumos_workflow`. The API has no unbind operation for these slots, so
`terraform destroy` drops the binding from state without calling the API and the
workflow stays bound. Each resource's page documents this.

Everything else in the V1 API is exposed read-only.

Not every V1 endpoint maps onto Terraform. The generated SDK covers all 61 V1
operations, but 13 of them are imperative actions or multipart uploads with no
declarative representation — triggering an app sync, completing or reassigning a
task, uploading a knowledge-context document. Those are reachable through the
SDK, not through a resource or data source.

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
      source  = "teamlumos/lumos"
      version = "1.0.0-alpha.1"
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
resource "lumos_app" "example" {
  app_class_id = "okta.com"
  auth = {
    # credentials vary by integration - discover them with the
    # lumos_integration_connection_schema data source
  }
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

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---

## SDK Created by Speakeasy

This SDK is generated and maintained using [Speakeasy](https://docs.speakeasyapi.dev/docs/using-speakeasy/client-sdks).
<!-- Start Summary [summary] -->
## Summary

Lumos: ALPHA: the 1.x line targets the Lumos V1 API and may introduce breaking changes without prior notice. Pin an exact version. For production, stay on 0.10.x (version = "~> 0.10").
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
      version = "1.0.0-alpha.1"
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

### Managed Resources

* [lumos_app](docs/resources/app.md)
* [lumos_app_approval_workflow](docs/resources/app_approval_workflow.md)
* [lumos_app_deprovisioning_workflow](docs/resources/app_deprovisioning_workflow.md)
* [lumos_app_provisioning_workflow](docs/resources/app_provisioning_workflow.md)
* [lumos_entitlement](docs/resources/entitlement.md)
* [lumos_entitlement_approval_workflow](docs/resources/entitlement_approval_workflow.md)
* [lumos_entitlement_deprovisioning_workflow](docs/resources/entitlement_deprovisioning_workflow.md)
* [lumos_entitlement_provisioning_workflow](docs/resources/entitlement_provisioning_workflow.md)
* [lumos_workflow](docs/resources/workflow.md)

### Data Sources

* [lumos_account_entitlements](docs/data-sources/account_entitlements.md)
* [lumos_app](docs/data-sources/app.md)
* [lumos_app_approval_workflow](docs/data-sources/app_approval_workflow.md)
* [lumos_app_deprovisioning_workflow](docs/data-sources/app_deprovisioning_workflow.md)
* [lumos_app_entitlement_app_store_catalog](docs/data-sources/app_entitlement_app_store_catalog.md)
* [lumos_app_entitlements](docs/data-sources/app_entitlements.md)
* [lumos_app_provisioning_workflow](docs/data-sources/app_provisioning_workflow.md)
* [lumos_app_store_catalog](docs/data-sources/app_store_catalog.md)
* [lumos_app_store_settings](docs/data-sources/app_store_settings.md)
* [lumos_apps](docs/data-sources/apps.md)
* [lumos_current_user](docs/data-sources/current_user.md)
* [lumos_custom_attribute_definitions](docs/data-sources/custom_attribute_definitions.md)
* [lumos_entitlement](docs/data-sources/entitlement.md)
* [lumos_entitlement_accounts](docs/data-sources/entitlement_accounts.md)
* [lumos_entitlement_app_store_settings](docs/data-sources/entitlement_app_store_settings.md)
* [lumos_entitlement_approval_workflow](docs/data-sources/entitlement_approval_workflow.md)
* [lumos_entitlement_deprovisioning_workflow](docs/data-sources/entitlement_deprovisioning_workflow.md)
* [lumos_entitlement_provisioning_workflow](docs/data-sources/entitlement_provisioning_workflow.md)
* [lumos_integration](docs/data-sources/integration.md)
* [lumos_integration_connection_schema](docs/data-sources/integration_connection_schema.md)
* [lumos_integrations](docs/data-sources/integrations.md)
* [lumos_knowledge_context](docs/data-sources/knowledge_context.md)
* [lumos_knowledge_contexts](docs/data-sources/knowledge_contexts.md)
* [lumos_task](docs/data-sources/task.md)
* [lumos_task_actions](docs/data-sources/task_actions.md)
* [lumos_tasks](docs/data-sources/tasks.md)
* [lumos_user](docs/data-sources/user.md)
* [lumos_user_accounts](docs/data-sources/user_accounts.md)
* [lumos_user_roles](docs/data-sources/user_roles.md)
* [lumos_users](docs/data-sources/users.md)
* [lumos_workflow](docs/data-sources/workflow.md)
* [lumos_workflow_steps](docs/data-sources/workflow_steps.md)
* [lumos_workflows](docs/data-sources/workflows.md)
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
