terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.13.1"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}