terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.11.2"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}