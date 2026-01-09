terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.10.0"
    }
  }
}

provider "lumos" {
  server_url = "..." # Required
}