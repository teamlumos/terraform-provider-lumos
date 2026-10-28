terraform {
  required_version = ">= 1.7.0"

  required_providers {
    lumos = {
      source = "teamlumos/lumos"
      # The prerelease under test. Ignored when dev_overrides points at a
      # local build.
      version = "1.0.0-beta.1"
    }
  }
}

provider "lumos" {
  server_url = var.server_url
}
