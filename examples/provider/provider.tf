terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.10.6"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}