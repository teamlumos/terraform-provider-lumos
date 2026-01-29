terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.12.3"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}