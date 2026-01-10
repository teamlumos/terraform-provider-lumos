terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "0.10.2"
    }
  }
}

provider "lumos" {
  server_url = "..." # Optional
}