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