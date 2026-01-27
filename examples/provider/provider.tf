terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.12.2"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}