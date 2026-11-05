terraform {
  required_version = ">= 1.7.0"

  required_providers {
    lumos = {
      source = "teamlumos/lumos"
      # The release under test. Ignored when dev_overrides points at a local
      # build.
      version = "1.0.0"
    }
  }
}

provider "lumos" {
  server_url = var.server_url
}
